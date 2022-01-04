package stream

import (
	"bytes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"net"

	"github.com/isayme/go-bufferpool"
	"github.com/pkg/errors"
)

// cipherInfo cipher definition
type cipherInfo struct {
	KeySize int
	IvSize  int

	genEncryptStream func(key, iv []byte) (cipher.Stream, error)
	genDecryptStream func(key, iv []byte) (cipher.Stream, error)
}

// Cipher cipher
type Cipher struct {
	enc cipher.Stream
	dec cipher.Stream

	buffer *bytes.Buffer

	key []byte

	info *cipherInfo
}

// NewCipher create cipher
func NewCipher(method string, key []byte) *Cipher {
	c := &Cipher{}
	c.key = key

	info, ok := cipherMethods[method]
	if !ok {
		panic(fmt.Errorf("unsupported method: %s", method))
	}

	c.info = info

	c.buffer = bytes.NewBuffer(nil)

	return c
}

// KeySize return key size
func (c *Cipher) KeySize() int {
	return c.info.KeySize
}

// getEncryptStream get encrypt stream
func (c *Cipher) getEncryptStream(iv []byte) (s cipher.Stream, err error) {
	_, err = io.ReadFull(rand.Reader, iv)
	if err != nil {
		return nil, err
	}

	s, err = c.info.genEncryptStream(c.key, iv)
	if err != nil {
		return nil, err
	}

	return s, nil
}

// getDecryptStream get decrypt stream
func (c *Cipher) getDecryptStream(iv []byte) (cipher.Stream, error) {
	return c.info.genDecryptStream(c.key, iv)
}

func (c *Cipher) decrypt(dst, src []byte) {
	c.dec.XORKeyStream(dst, src)
}

func (c *Cipher) encrypt(dst, src []byte) {
	c.enc.XORKeyStream(dst, src)
}

// Read read from client
func (c *Cipher) Read(conn net.Conn, p []byte) (n int, err error) {
	if c.dec == nil {
		iv := bufferpool.Get(c.info.IvSize)
		defer bufferpool.Put(iv)

		if _, err = io.ReadFull(conn, iv); err != nil {
			return 0, err
		}

		s, err := c.getDecryptStream(iv)
		if err != nil {
			return 0, err
		}

		c.dec = s
	}

	n, err = conn.Read(p)
	c.decrypt(p, p[0:n])
	return n, err
}

// Write write to client
func (c *Cipher) Write(conn net.Conn, p []byte) (n int, err error) {
	if c.enc == nil {
		iv := bufferpool.Get(c.info.IvSize)
		defer bufferpool.Put(iv)

		c.enc, err = c.getEncryptStream(iv)
		if err != nil {
			return 0, err
		}

		nw, err := conn.Write(iv)
		if err != nil {
			return 0, errors.Wrap(err, "iv write")
		}
		if nw != len(iv) {
			return 0, errors.New("iv short write")
		}
	}

	c.encrypt(p, p)
	return conn.Write(p)
}

package stream

import (
	"bytes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"net"

	"github.com/isayme/go-shadowsocks/shadowsocks/bufferpool"
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
	Method string

	Enc cipher.Stream
	Dec cipher.Stream

	buffer *bytes.Buffer

	key []byte

	Info *cipherInfo
}

// NewCipher create cipher
func NewCipher(method string, key []byte) *Cipher {
	c := &Cipher{}
	c.Method = method
	c.key = key

	Info, ok := cipherMethods[method]
	if !ok {
		panic(fmt.Errorf("unsupported method: %s", method))
	}

	c.Info = Info

	c.buffer = bytes.NewBuffer(nil)

	return c
}

// KeySize return key size
func (c *Cipher) KeySize() int {
	return c.Info.KeySize
}

// getEncryptStream get encrypt stream
func (c *Cipher) getEncryptStream(iv []byte) (s cipher.Stream, err error) {
	_, err = io.ReadFull(rand.Reader, iv)
	if err != nil {
		return nil, err
	}

	s, err = c.Info.genEncryptStream(c.key, iv)
	if err != nil {
		return nil, err
	}

	return s, nil
}

// getDecryptStream get decrypt stream
func (c *Cipher) getDecryptStream(iv []byte) (cipher.Stream, error) {
	return c.Info.genDecryptStream(c.key, iv)
}

func (c *Cipher) decrypt(dst, src []byte) {
	c.Dec.XORKeyStream(dst, src)
}

func (c *Cipher) encrypt(dst, src []byte) {
	c.Enc.XORKeyStream(dst, src)
}

// Read read from client
func (c *Cipher) Read(conn net.Conn, p []byte) (n int, err error) {
	if c.Dec == nil {
		iv := bufferpool.Get(c.Info.IvSize)
		defer bufferpool.Put(iv)

		if _, err = io.ReadFull(conn, iv); err != nil {
			return 0, err
		}

		s, err := c.getDecryptStream(iv)
		if err != nil {
			return 0, err
		}

		c.Dec = s
	}

	n, err = conn.Read(p)
	c.decrypt(p, p[0:n])
	return n, err
}

// Write write to client
func (c *Cipher) Write(conn net.Conn, p []byte) (n int, err error) {
	if c.Enc == nil {
		iv := bufferpool.Get(c.Info.IvSize)
		defer bufferpool.Put(iv)

		c.Enc, err = c.getEncryptStream(iv)
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

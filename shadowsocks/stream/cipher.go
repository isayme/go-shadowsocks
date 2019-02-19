package stream

import (
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"net"

	"github.com/isayme/go-shadowsocks/shadowsocks/bufferpool"
)

// cipherInfo cipher definition
type cipherInfo struct {
	KeyLen int
	IvLen  int

	genEncryptStream func(key, iv []byte) (cipher.Stream, error)
	genDecryptStream func(key, iv []byte) (cipher.Stream, error)
}

// Cipher cipher
type Cipher struct {
	Method   string
	Password string

	Conn net.Conn

	Enc cipher.Stream
	Dec cipher.Stream

	key   []byte
	nonce []byte

	Info *cipherInfo
}

// GetEncryptStream get encrypt stream
func (c *Cipher) GetEncryptStream(iv []byte) (s cipher.Stream, err error) {
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

// GetDecryptStream get decrypt stream
func (c *Cipher) GetDecryptStream(iv []byte) (cipher.Stream, error) {
	return c.Info.genDecryptStream(c.key, iv)
}

func (c *Cipher) decrypt(dst, src []byte) {
	c.Dec.XORKeyStream(dst, src)
}

func (c *Cipher) encrypt(dst, src []byte) {
	c.Enc.XORKeyStream(dst, src)
}

// Read read from client
func (c *Cipher) Read(p []byte) (n int, err error) {
	if c.Dec == nil {
		iv := bufferpool.Get(c.Info.IvLen)
		defer bufferpool.Put(iv)

		if _, err = io.ReadFull(c.Conn, iv); err != nil {
			return 0, err
		}

		s, err := c.GetDecryptStream(iv)
		if err != nil {
			return 0, err
		}

		c.Dec = s
	}

	n, err = c.Conn.Read(p)
	c.decrypt(p, p[0:n])
	return n, err
}

// Write write to client
func (c *Cipher) Write(p []byte) (n int, err error) {
	if c.Enc == nil {
		iv := bufferpool.Get(c.Info.IvLen)
		defer bufferpool.Put(iv)

		c.Enc, err = c.GetEncryptStream(iv)
		if err != nil {
			return 0, err
		}

		_, err = c.Conn.Write(iv)
		if err != nil {
			return 0, err
		}
	}

	c.encrypt(p, p)
	return c.Conn.Write(p)
}

// NewCipher create cipher
func NewCipher(method string, password string, conn net.Conn) (*Cipher, error) {
	c := &Cipher{}
	c.Method = method
	c.Password = password
	c.Conn = conn

	Info, ok := cipherMethods[method]
	if !ok {
		return nil, fmt.Errorf("unsupported method: %s", method)
	}

	c.Info = Info
	c.key = generateKey(c.Password, c.Info.KeyLen)

	return c, nil
}

package cipher

import (
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
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

	Enc cipher.Stream
	Dec cipher.Stream

	key []byte

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

// Clone clone to ignore key generate
func (c *Cipher) Clone() *Cipher {
	nc := *c
	nc.Dec = nil
	nc.Dec = nil
	return &nc
}

// NewCipher create cipher
func NewCipher(method string, password string) (*Cipher, error) {
	c := &Cipher{}
	c.Method = method
	c.Password = password

	Info, ok := cipherMethods[method]
	if !ok {
		return nil, fmt.Errorf("unsupported method: %s", method)
	}

	c.Info = Info
	c.key = generateKey(c.Password, c.Info.KeyLen)

	return c, nil
}

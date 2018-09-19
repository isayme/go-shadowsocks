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

	*cipherInfo
}

// GetEncryptStream get encrypt stream
func (c Cipher) GetEncryptStream() (iv []byte, s cipher.Stream, err error) {
	key := generateKey(c.Password, c.KeyLen)

	iv = make([]byte, c.IvLen)
	_, err = io.ReadFull(rand.Reader, iv)
	if err != nil {
		return nil, nil, err
	}

	s, err = c.genEncryptStream(key, iv)
	if err != nil {
		return nil, nil, err
	}

	return iv, s, nil
}

// GetDecryptStream get decrypt stream
func (c Cipher) GetDecryptStream(iv []byte) (cipher.Stream, error) {
	key := generateKey(c.Password, c.KeyLen)

	return c.genDecryptStream(key, iv)
}

// NewCipher create cipher
func NewCipher(method string, password string) (*Cipher, error) {
	c := &Cipher{}
	c.Method = method
	c.Password = password

	info, ok := cipherMethods[method]
	if !ok {
		return nil, fmt.Errorf("unsupported method: %s", method)
	}

	c.cipherInfo = info

	return c, nil
}

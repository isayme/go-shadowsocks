package stream

import (
	"crypto/cipher"
	"fmt"
	"io"
)

// NewCipher create cipher
func NewCipher(rw io.ReadWriter, method string, key []byte) *cipherInfo {
	info, ok := cipherMethods[method]
	if !ok {
		panic(fmt.Errorf("unsupported method: %s", method))
	}

	c := info
	c.Writer = NewWriter(rw, key, info.ivSize, info.newWriter)
	c.Reader = NewReader(rw, key, info.ivSize, info.newReader)

	return &c
}

// cipherInfo cipher definition
type cipherInfo struct {
	keySize int
	ivSize  int

	newWriter func(key, iv []byte) (cipher.Stream, error)
	newReader func(key, iv []byte) (cipher.Stream, error)

	io.Reader
	io.Writer
}

func newCipherInfo(keySize, ivSize int, newReader, newWriter func(key, iv []byte) (cipher.Stream, error)) cipherInfo {
	return cipherInfo{
		keySize:   keySize,
		ivSize:    ivSize,
		newReader: newReader,
		newWriter: newWriter,
	}
}

// KeySize return key size
func (c *cipherInfo) KeySize() int {
	return c.keySize
}

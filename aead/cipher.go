package aead

import (
	"crypto/cipher"
	"fmt"
	"io"
)

// NewCipher create aead cipher
func NewCipher(rw io.ReadWriter, method string, key []byte) *cipherInfo {
	info, ok := cipherMethods[method]
	if !ok {
		panic(fmt.Errorf("unsupported method: %s", method))
	}

	c := info
	c.Reader = NewReader(rw, key, info.newReader)
	c.Writer = NewWriter(rw, key, info.newWriter)

	return &c
}

// cipherInfo cipher definition
type cipherInfo struct {
	keySize int

	newWriter func([]byte) (cipher.AEAD, error)
	newReader func([]byte) (cipher.AEAD, error)

	io.Reader
	io.Writer
}

func newCipherInfo(keySize int, newReader, newWriter func([]byte) (cipher.AEAD, error)) cipherInfo {
	return cipherInfo{
		keySize:   keySize,
		newReader: newReader,
		newWriter: newWriter,
	}
}

func (ci *cipherInfo) KeySize() int {
	return ci.keySize
}

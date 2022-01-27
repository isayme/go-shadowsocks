package aead

import (
	"crypto/aes"
	"crypto/cipher"
)

func newAesGcmCipher(key []byte) (cipher.AEAD, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	return cipher.NewGCM(block)
}

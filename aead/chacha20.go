package aead

import (
	"crypto/cipher"

	"golang.org/x/crypto/chacha20poly1305"
)

func newChacha20Poly1305Cipher(key []byte) (cipher.AEAD, error) {
	return chacha20poly1305.New(key)
}

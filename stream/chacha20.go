package stream

import (
	"crypto/cipher"

	"github.com/aead/chacha20"
)

func newChaCha20Writer(key, iv []byte) (cipher.Stream, error) {
	return chacha20.NewCipher(iv, key)
}

var newChaCha20Reader = newChaCha20Writer

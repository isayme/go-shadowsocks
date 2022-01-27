package stream

import (
	"crypto/cipher"

	"golang.org/x/crypto/blowfish"
)

func newBlowfishWriter(key, iv []byte) (cipher.Stream, error) {
	block, err := blowfish.NewCipher(key)
	if err != nil {
		return nil, err
	}

	return cipher.NewCFBEncrypter(block, iv), nil
}

func newBlowfishReader(key, iv []byte) (cipher.Stream, error) {
	block, err := blowfish.NewCipher(key)
	if err != nil {
		return nil, err
	}

	return cipher.NewCFBDecrypter(block, iv), nil
}

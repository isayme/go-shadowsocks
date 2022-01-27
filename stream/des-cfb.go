package stream

import (
	"crypto/cipher"
	"crypto/des"
)

func newDESCFBWriter(key, iv []byte) (cipher.Stream, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}

	return cipher.NewCFBEncrypter(block, iv), nil
}

func newDESCFBReader(key, iv []byte) (cipher.Stream, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}

	return cipher.NewCFBDecrypter(block, iv), nil
}

package stream

import (
	"crypto/cipher"
	"crypto/rc4"

	"github.com/isayme/go-shadowsocks/util"
)

func newRC4MD5Writer(key, iv []byte) (cipher.Stream, error) {
	rc4key := util.MD5(key, iv)

	return rc4.NewCipher(rc4key)
}

var newRC4MD5Reader = newRC4MD5Writer

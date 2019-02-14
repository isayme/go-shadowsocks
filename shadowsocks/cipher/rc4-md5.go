package cipher

import (
	"crypto/cipher"
	"crypto/md5"
	"crypto/rc4"
)

func newRC4MD5Stream(key, iv []byte) (cipher.Stream, error) {
	h := md5.New()
	h.Write(key)
	h.Write(iv)
	rc4key := h.Sum(nil)

	return rc4.NewCipher(rc4key)
}

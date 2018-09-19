package util

import "crypto/md5"

func MD5(p []byte) []byte {
	h := md5.New()
	h.Write(p)
	return h.Sum(nil)
}

package cipher

import (
	"crypto/md5"
	"math"

	"github.com/isayme/go-shadowsocks/shadowsocks/util"
)

func generateKey(password string, keyLen int) []byte {
	count := int(math.Ceil(float64(keyLen) / float64(md5.Size)))

	r := make([]byte, count*md5.Size)

	copy(r, util.MD5([]byte(password)))

	d := make([]byte, md5.Size+len(password))
	start := 0
	for i := 1; i < count; i++ {
		start += md5.Size
		copy(d[:md5.Size], r[start-md5.Size:start])
		copy(d[md5.Size:], password)
		copy(r[start:start+md5.Size], util.MD5(d))
	}

	return r[:keyLen]
}

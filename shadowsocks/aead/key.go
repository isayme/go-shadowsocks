package aead

import (
	"crypto/md5"
	"crypto/sha1"
	"io"
	"math"

	"golang.org/x/crypto/hkdf"

	"github.com/isayme/go-shadowsocks/shadowsocks/bufferpool"
	"github.com/isayme/go-shadowsocks/shadowsocks/util"
)

var hkdfInfo = []byte("ss-subkey")

func generateKey(password string, keyLen int) []byte {
	count := int(math.Ceil(float64(keyLen) / float64(md5.Size)))

	r := bufferpool.Get(count * md5.Size)
	defer bufferpool.Put(r)

	copy(r, util.MD5([]byte(password)))

	d := bufferpool.Get(md5.Size + len(password))
	defer bufferpool.Put(d)

	start := 0
	for i := 1; i < count; i++ {
		start += md5.Size
		copy(d[:md5.Size], r[start-md5.Size:start])
		copy(d[md5.Size:], password)
		copy(r[start:start+md5.Size], util.MD5(d))
	}

	key := make([]byte, keyLen)
	copy(key, r[:keyLen])

	return key
}

func hkdfSHA1(secret, salt, info, subkey []byte) {
	r := hkdf.New(sha1.New, secret, salt, info)
	if _, err := io.ReadFull(r, subkey); err != nil {
		panic(err) // should never happen
	}
}

func increment(b []byte) {
	for i := range b {
		b[i]++
		if b[i] != 0 {
			return
		}
	}
}

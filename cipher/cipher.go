package cipher

import (
	"io"

	"github.com/isayme/go-shadowsocks/aead"
	"github.com/isayme/go-shadowsocks/stream"
	"github.com/isayme/go-shadowsocks/util"
)

// Cipher cipher interface
type Cipher interface {
	Read(p []byte) (n int, err error)
	Write(p []byte) (n int, err error)
	KeySize() int
}

// NewCipher create cipher
func NewCipher(rw io.ReadWriter, method string, key []byte) Cipher {
	switch method {
	case "aes-128-gcm", "aes-192-gcm", "aes-256-gcm", "chacha20-ietf-poly1305":
		return aead.NewCipher(rw, method, key)
	default:
		return stream.NewCipher(rw, method, key)
	}
}

// NewKey create key from method/password
func NewKey(method, password string) []byte {
	cipher := NewCipher(nil, method, nil)

	keySize := cipher.KeySize()
	return util.KDF(password, keySize)
}

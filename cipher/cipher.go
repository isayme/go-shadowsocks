package cipher

import (
	"net"

	"github.com/isayme/go-shadowsocks/aead"
	"github.com/isayme/go-shadowsocks/stream"
	"github.com/isayme/go-shadowsocks/util"
)

// Cipher cipher interface
type Cipher interface {
	Read(conn net.Conn, b []byte) (n int, err error)
	Write(conn net.Conn, b []byte) (n int, err error)
	KeySize() int
}

// NewCipher create cipher
func NewCipher(method string, key []byte) Cipher {
	switch method {
	case "aes-128-gcm", "aes-192-gcm", "aes-256-gcm", "chacha20-ietf-poly1305":
		return aead.NewCipher(method, key)
	default:
		return stream.NewCipher(method, key)
	}
}

// NewKey create key from method/password
func NewKey(method, password string) []byte {
	cipher := NewCipher(method, nil)

	keySize := cipher.KeySize()
	return util.KDF(password, keySize)
}

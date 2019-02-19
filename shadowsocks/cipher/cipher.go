package cipher

import (
	"net"

	"github.com/isayme/go-shadowsocks/shadowsocks/aead"
	"github.com/isayme/go-shadowsocks/shadowsocks/stream"
)

// Cipher cipher interface
type Cipher interface {
	Read(p []byte) (n int, err error)
	Write(p []byte) (n int, err error)
}

// NewCipher create cipher
func NewCipher(method, password string, conn net.Conn) (Cipher, error) {
	switch method {
	case "aes-128-gcm", "aes-192-gcm", "aes-256-gcm":
		return aead.NewCipher(method, password, conn)
	default:
		return stream.NewCipher(method, password, conn)
	}
}

package cipher

import (
	"net"

	"github.com/isayme/go-shadowsocks/shadowsocks/aead"
	"github.com/isayme/go-shadowsocks/shadowsocks/stream"
	"github.com/isayme/go-shadowsocks/shadowsocks/util"
)

// Cipher cipher interface
type Cipher interface {
	Read(p []byte) (n int, err error)
	Write(p []byte) (n int, err error)
	KeySize() int
	Init(key []byte, conn net.Conn)
}

// NewCipher create cipher
func NewCipher(method string) Cipher {
	switch method {
	case "aes-128-gcm", "aes-192-gcm", "aes-256-gcm", "chacha20-ietf-poly1305":
		return aead.NewCipher(method)
	default:
		return stream.NewCipher(method)
	}
}

// NewKey create key from method/password
func NewKey(method, password string) []byte {
	cipher := NewCipher(method)

	keySize := cipher.KeySize()
	return util.KDF(password, keySize)
}

type CipherConn struct {
	net.Conn
	cipher Cipher
}

// NewCipherConn create remote instance
func NewCipherConn(conn net.Conn, c Cipher) *CipherConn {
	return &CipherConn{
		Conn:   conn,
		cipher: c,
	}
}

// Read read from remote
func (conn *CipherConn) Read(p []byte) (n int, err error) {
	return conn.cipher.Read(p)
}

// Write write to conn
func (conn *CipherConn) Write(p []byte) (n int, err error) {
	return conn.cipher.Write(p)
}

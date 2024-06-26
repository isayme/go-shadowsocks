package ss

import (
	"encoding/binary"
	"io"
	"net"
	"strconv"

	"github.com/isayme/go-bufferpool"
	"github.com/isayme/go-logger"
	"github.com/isayme/go-shadowsocks/cipher"
	"github.com/isayme/go-shadowsocks/util"
	"github.com/pkg/errors"
)

type Connection struct {
	net.Conn
	cipher cipher.Cipher
}

func NewConnection(conn net.Conn, method string, key []byte) Connection {
	cipher := cipher.NewCipher(conn, method, key)

	return Connection{
		Conn:   conn,
		cipher: cipher,
	}
}

func (c Connection) Read(b []byte) (n int, err error) {
	return c.cipher.Read(b)
}

func (c Connection) Write(b []byte) (n int, err error) {
	return c.cipher.Write(b)
}

/**
 * read address to be proxyed.
 * used for connection from shadowsocks client
 */
func (c Connection) readAddress() (string, error) {
	data := bufferpool.Get(256)
	defer bufferpool.Put(data)

	if _, err := io.ReadFull(c, data[:1]); err != nil {
		return "", errors.Wrap(err, "read type")
	}

	typ := data[0]
	logger.Debugf("address type: %02x", typ)

	var host string
	switch typ {
	case util.AddressTypeIPV4:
		if _, err := io.ReadFull(c, data[:net.IPv4len]); err != nil {
			return "", errors.Wrap(err, "read ipv4")
		}
		host = net.IP(data[:net.IPv4len]).String()
	case util.AddressTypeDomain:
		if _, err := io.ReadFull(c, data[:1]); err != nil {
			return "", errors.Wrap(err, "read domain length")
		}
		domainLen := int(data[0])

		if _, err := io.ReadFull(c, data[:domainLen]); err != nil {
			return "", errors.Wrap(err, "read domain")
		}
		host = string(data[:domainLen])
	case util.AddressTypeIPV6:
		if _, err := io.ReadFull(c, data[:net.IPv6len]); err != nil {
			return "", errors.Wrap(err, "read ipv6")
		}
		host = net.IP(data[:net.IPv6len]).String()
	default:
		return "", errors.Errorf("invalid address type: %02x", typ)
	}

	if _, err := io.ReadFull(c, data[:2]); err != nil {
		return "", errors.Wrap(err, "read port")
	}

	port := binary.BigEndian.Uint16(data)

	return net.JoinHostPort(host, strconv.Itoa(int(port))), nil
}

package connection

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/isayme/go-shadowsocks/shadowsocks/bufferpool"
	"github.com/isayme/go-shadowsocks/shadowsocks/cipher"
	"github.com/isayme/go-shadowsocks/shadowsocks/logger"
	"github.com/pkg/errors"
)

// Client client from sslocal
type Client struct {
	Conn net.Conn

	Cipher cipher.Cipher
}

// NewClient create client instance
func NewClient(conn net.Conn, c cipher.Cipher) (*Client, error) {
	client := &Client{
		Conn:   conn,
		Cipher: c,
	}

	return client, nil
}

// Close close connection
func (client Client) Close() error {
	return client.Conn.Close()
}

// Read read from client
func (client *Client) Read(p []byte) (n int, err error) {
	return client.Cipher.Read(p)
}

// Write write to client
func (client *Client) Write(p []byte) (n int, err error) {
	return client.Cipher.Write(p)
}

// SetReadTimeout set read timeout
func (client *Client) SetReadTimeout(timeout int) {
	if timeout <= 0 {
		client.Conn.SetReadDeadline(time.Time{})
	} else {
		client.Conn.SetReadDeadline(time.Now().Add(time.Second * time.Duration(timeout)))
	}
}

// ReadAddress read & parse remote address
func (client *Client) ReadAddress(timeout int) (string, error) {
	client.SetReadTimeout(timeout)
	defer client.SetReadTimeout(0)

	data := bufferpool.Get(256)
	defer bufferpool.Put(data)

	if _, err := io.ReadFull(client, data[:1]); err != nil {
		return "", errors.Wrap(err, "read type")
	}

	typ := AddressType(data[0])
	if !typ.Valid() {
		return "", errors.Errorf("invalid address type: %02x", typ)
	}
	logger.Debugf("address type: %s", typ)

	var host string
	switch typ {
	case AddressTypeIPV4:
		if _, err := io.ReadFull(client, data[:net.IPv4len]); err != nil {
			return "", errors.Wrap(err, "read ipv4")
		}
		host = net.IP(data[:net.IPv4len]).String()
	case AddressTypeDomain:
		if _, err := io.ReadFull(client, data[:1]); err != nil {
			return "", errors.Wrap(err, "read domain length")
		}
		domainLen := int(data[0])

		if _, err := io.ReadFull(client, data[:domainLen]); err != nil {
			return "", errors.Wrap(err, "read domain")
		}
		host = string(data[:domainLen])
	case AddressTypeIPV6:
		if _, err := io.ReadFull(client, data[:net.IPv6len]); err != nil {
			return "", errors.Wrap(err, "read ipv6")
		}
		host = net.IP(data[:net.IPv6len]).String()
	}

	if _, err := io.ReadFull(client, data[:2]); err != nil {
		return "", errors.Wrap(err, "read port")
	}

	var port uint16
	port = binary.BigEndian.Uint16(data)

	return fmt.Sprintf("%s:%d", host, port), nil
}

package connection

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"

	"github.com/isayme/go-shadowsocks/shadowsocks/cipher"
	"github.com/isayme/go-shadowsocks/shadowsocks/logger"
	"github.com/pkg/errors"
)

// Client client from sslocal
type Client struct {
	Conn net.Conn

	cipher.Cipher
}

// NewClient create client instance
func NewClient(conn net.Conn, c cipher.Cipher) (*Client, error) {
	client := &Client{
		Conn:   conn,
		Cipher: c,
	}

	return client, nil
}

func (client Client) decrypt(dst, src []byte) {
	client.Dec.XORKeyStream(dst, src)
}

func (client Client) encrypt(dst, src []byte) {
	client.Enc.XORKeyStream(dst, src)
}

// Close close connection
func (client Client) Close() error {
	return client.Conn.Close()
}

// Read read from client
func (client *Client) Read(p []byte) (n int, err error) {
	if client.Dec == nil {
		iv := make([]byte, client.IvLen)
		if _, err = io.ReadFull(client.Conn, iv); err != nil {
			return 0, err
		}

		s, err := client.GetDecryptStream(iv)
		if err != nil {
			return 0, err
		}

		client.Dec = s
	}

	data := make([]byte, len(p))

	n, err = client.Conn.Read(data)
	client.decrypt(p, data[0:n])
	return n, err
}

// Write write to client
func (client *Client) Write(p []byte) (n int, err error) {
	var iv []byte

	if client.Enc == nil {
		iv, client.Enc, err = client.GetEncryptStream()
		if err != nil {
			return 0, err
		}
	}

	data := make([]byte, len(iv)+len(p))
	copy(data, iv)

	client.encrypt(data[len(iv):], p)
	n, err = client.Conn.Write(data)
	return n - len(iv), err
}

// ReadAddress read & parse remote address
func (client *Client) ReadAddress() (string, error) {
	var data []byte

	data = make([]byte, 1)
	if _, err := io.ReadFull(client, data); err != nil {
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
		data = make([]byte, net.IPv4len)
		if _, err := io.ReadFull(client, data); err != nil {
			return "", errors.Wrap(err, "read ipv4")
		}
		host = net.IP(data).String()
	case AddressTypeDomain:
		data = make([]byte, 1)
		if _, err := io.ReadFull(client, data); err != nil {
			return "", errors.Wrap(err, "read domain length")
		}
		domainLen := int(data[0])

		data = make([]byte, domainLen)
		if _, err := io.ReadFull(client, data); err != nil {
			return "", errors.Wrap(err, "read domain")
		}
		host = string(data)
	case AddressTypeIPV6:
		data = make([]byte, net.IPv6len)
		if _, err := io.ReadFull(client, data); err != nil {
			return "", errors.Wrap(err, "read ipv6")
		}
		host = net.IP(data).String()
	}

	data = make([]byte, 2)
	if _, err := io.ReadFull(client, data); err != nil {
		return "", errors.Wrap(err, "read port")
	}

	var port uint16
	port = binary.BigEndian.Uint16(data)

	return fmt.Sprintf("%s:%d", host, port), nil
}

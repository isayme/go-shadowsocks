package connection

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"

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

	n, err = client.Conn.Read(p)
	client.decrypt(p, p[0:n])
	return n, err
}

// Write write to client
func (client *Client) Write(p []byte) (n int, err error) {
	if client.Enc == nil {
		var iv []byte
		iv, client.Enc, err = client.GetEncryptStream()
		if err != nil {
			return 0, err
		}

		_, err = client.Conn.Write(iv)
		if err != nil {
			return 0, err
		}
	}

	client.encrypt(p, p)
	return client.Conn.Write(p)
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

	data := make([]byte, 256)

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

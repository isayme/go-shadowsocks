package socks5

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"

	"github.com/isayme/go-shadowsocks/shadowsocks/util"
	"github.com/pkg/errors"
)

// Version socks5 version
const Version = 5

// method
const (
	MethodNone = 0 // NO AUTHENTICATION REQUIRED
)

// request cmd
const (
	CmdConnect      = 0x01
	CmdUDPAssociate = 0x03
)

// address type
const (
	AddressTypeIPV4   = util.AddressTypeIPV4
	AddressTypeDomain = util.AddressTypeDomain
	AddressTypeIPV6   = util.AddressTypeIPV6
)

type Request struct {
	client net.Conn

	cmd     byte
	atyp    byte
	addr    string
	RawAddr []byte
}

func NewRequest(client net.Conn) (*Request, error) {
	request := &Request{
		client: client,
	}

	err := request.negotiate()
	if err != nil {
		return nil, err
	}

	return request, nil
}

func (r *Request) negotiate() error {
	buf := make([]byte, 256)

	// version
	_, err := io.ReadFull(r.client, buf[:1])
	if err != nil {
		return errors.Errorf("read version fail: %s", err)
	}

	if buf[0] != Version {
		return errors.New("not socks5 protocol")
	}

	// methods
	_, err = io.ReadFull(r.client, buf[:1])
	if err != nil {
		return errors.Errorf("read nmethods fail: %s", err)
	}
	nMethods := buf[0]
	if nMethods < 1 {
		return errors.Errorf("nmethods not valid: %d", nMethods)
	}

	_, err = io.ReadFull(r.client, buf[:nMethods])
	if err != nil {
		return errors.Errorf("read nmethods fail: %s", err)
	}

	_, err = r.client.Write([]byte{Version, MethodNone})
	if err != nil {
		return errors.Errorf("write accepet method fail: %s", err)
	}

	_, err = io.ReadFull(r.client, buf[:4])
	if err != nil {
		return errors.Errorf("read adrress fail: %s", err)
	}
	r.cmd = buf[1]
	r.atyp = buf[3]

	var reply = []byte{Version, 0, 0, r.atyp}

	switch r.cmd {
	case CmdConnect:
	default:
		return errors.Errorf("not support cmd: %d", r.cmd)
	}

	switch r.atyp {
	case AddressTypeDomain:
		_, err = io.ReadFull(r.client, buf[:1])
		if err != nil {
			return errors.Errorf("read adrress fail: %s", err)
		}
		domainLen := buf[0]
		reply = append(reply, buf[0])

		_, err = io.ReadFull(r.client, buf[:domainLen+2])
		if err != nil {
			return errors.Errorf("read adrress fail: %s", err)
		}
		reply = append(reply, buf[:domainLen+2]...)

		domain := string(buf[:domainLen])
		port := binary.BigEndian.Uint16(buf[domainLen : domainLen+2])
		r.addr = fmt.Sprintf("%s:%d", domain, port)
	case AddressTypeIPV4:
		_, err = io.ReadFull(r.client, buf[:6])
		if err != nil {
			return errors.Errorf("read adrress fail: %s", err)
		}

		reply = append(reply, buf[:6]...)

		ip := net.IPv4(buf[0], buf[1], buf[2], buf[3])
		port := binary.BigEndian.Uint16(buf[4:6])

		r.addr = fmt.Sprintf("%s:%d", ip.String(), port)
	// case AddressTypeIPV6:
	default:
		return errors.Errorf("not support adrress type: %d", r.atyp)
	}

	r.RawAddr = reply[3:]
	_, err = r.client.Write(reply)
	if err != nil {
		return errors.Errorf("reply request fail: %s", err)
	}

	return nil
}

func (r *Request) RemoteAddress() string {
	return r.addr
}

package ss

import (
	"net"
	"strconv"
	"time"

	"github.com/isayme/go-logger"
	"github.com/isayme/go-shadowsocks/cipher"
	"github.com/isayme/go-shadowsocks/socks5"
	"github.com/isayme/go-shadowsocks/util"
)

type Client struct {
	method  string
	key     []byte
	timeout time.Duration

	server string
	port   int
}

func NewClient(server string, port int, method string, password string, timeout time.Duration) *Client {
	key := cipher.NewKey(method, password)

	return &Client{
		method:  method,
		key:     key,
		timeout: timeout,

		server: server,
		port:   port,
	}
}

func (c *Client) AcceptAndHandle(conn net.Conn) {
	logger.Debugf("new connection from: %s", conn.RemoteAddr().String())
	client := util.NewTimeoutConn(conn, c.timeout)
	defer client.Close()

	request, err := socks5.NewRequest(client)
	if err != nil {
		logger.Errorf("NewRequest fail, err: %s", err)
		return
	}

	address := net.JoinHostPort(c.server, strconv.Itoa(c.port))
	ssconn, err := net.DialTimeout("tcp", address, 5*time.Second)
	if err != nil {
		logger.Errorf("dial ssserver fail, err: %s", err)
		return
	}

	remote := NewConnection(ssconn, c.method, c.key)
	defer remote.Close()

	_, err = remote.Write(request.RawAddr)
	if err != nil {
		logger.Errorf("dial ssserver fail, err: %s", err)
		return
	}

	util.Proxy(client, remote)

	logger.Debugf("connection '%s => %s' closed", conn.RemoteAddr().String(), request.RemoteAddress())
}

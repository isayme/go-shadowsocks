package ss

import (
	"net"
	"time"

	"github.com/isayme/go-logger"
	"github.com/isayme/go-shadowsocks/cipher"
	"github.com/isayme/go-shadowsocks/util"
)

type Server struct {
	method  string
	key     []byte
	timeout time.Duration
}

func NewServer(method string, password string, timeout time.Duration) *Server {
	key := cipher.NewKey(method, password)

	return &Server{
		method:  method,
		key:     key,
		timeout: timeout,
	}
}

func (s *Server) AcceptAndHandle(conn net.Conn) {
	logger.Debugf("new connection from: %s", conn.RemoteAddr().String())
	tcpConn, _ := conn.(*net.TCPConn)
	conn = util.NewTimeoutConn(conn, s.timeout)
	ssconn := NewConnection(conn, s.method, s.key)
	defer ssconn.Close()

	address, err := ssconn.readAddress()
	if err != nil {
		logger.Errorw("read address fail", "err", err, "remoteAddr", conn.RemoteAddr().String())
		return
	}

	logger.Infof("connecting remote '%s'", address)
	// dial with timeout
	remote, err := net.DialTimeout("tcp", address, time.Second*5)
	if err != nil {
		logger.Warnf("dial remote '%s' failed, err: %+v", address, err)
		return
	}
	tcpRemoteConn, _ := remote.(*net.TCPConn)
	remote = util.NewTimeoutConn(remote, s.timeout)
	defer remote.Close()
	logger.Debugf("connect remote '%s' success", address)

	util.Proxy(ssconn, tcpConn, remote, tcpRemoteConn)

	logger.Debugf("connection '%s => %s' closed", conn.RemoteAddr().String(), address)
}

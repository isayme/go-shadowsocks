package connection

import (
	"net"
	"time"
)

// Remote remote connect
type Remote struct {
	Conn    net.Conn
	Address string
}

// NewRemote create remote client
func NewRemote(conn net.Conn, address string) *Remote {
	return &Remote{
		Conn:    conn,
		Address: address,
	}
}

// SetReadTimeout set read timeout
func (remote *Remote) SetReadTimeout(timeout int) {
	if timeout <= 0 {
		remote.Conn.SetReadDeadline(time.Time{})
	} else {
		remote.Conn.SetReadDeadline(time.Now().Add(time.Second * time.Duration(timeout)))
	}
}

// Read read from remote
func (remote *Remote) Read(p []byte) (n int, err error) {
	return remote.Conn.Read(p)
}

// Write write to remote
func (remote *Remote) Write(p []byte) (n int, err error) {
	return remote.Conn.Write(p)
}

// Close close connection
func (remote *Remote) Close() error {
	return remote.Conn.Close()
}

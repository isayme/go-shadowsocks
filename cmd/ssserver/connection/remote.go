package connection

import "net"

// Remote remote connect
type Remote struct {
	Conn net.Conn
}

// NewRemote create remote client
func NewRemote(conn net.Conn) *Remote {
	return &Remote{
		Conn: conn,
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

package connection

import (
	"io"
	"os"

	"github.com/isayme/go-shadowsocks/shadowsocks/bufferpool"

	"github.com/isayme/go-shadowsocks/shadowsocks/logger"
)

// Connection connection from client
type Connection struct {
	Remote  *Remote
	Client  *Client
	Timeout int
}

// Serve start exchange data between remote & client
func (c Connection) Serve() {
	// any of remote/client closed, the other one should close with quiet
	closed := false

	go func() {
		_, err := copyBuffer(c.Remote, c.Client, c.Timeout)
		if err != nil && !closed {
			logger.Errorf("io.Copy from client to remote fail, err: %#v", err)
		}
		closed = true
		logger.Debug("client read end")
		c.Remote.Close()
	}()

	_, err := copyBuffer(c.Client, c.Remote, c.Timeout)
	if err != nil && !closed {
		logger.Errorf("io.Copy from remote to client fail, err: %#v", err)
	}
	closed = true
	logger.Debug("remote read end")
	c.Client.Close()
}

var bufSize = os.Getpagesize()

// TimeoutReader reader with read timeout
type TimeoutReader interface {
	Read(p []byte) (n int, err error)
	SetReadTimeout(timeout int)
}

// code from io.CopyBuffer, paste here to enable read timeout
func copyBuffer(dst io.Writer, src TimeoutReader, timeout int) (written int64, err error) {
	buf := bufferpool.DefaultPool.Get()
	defer bufferpool.DefaultPool.Put(buf)

	for {
		src.SetReadTimeout(timeout)
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
	}
	return written, err
}

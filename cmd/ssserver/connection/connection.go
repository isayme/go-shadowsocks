package connection

import (
	"io"
	"os"
	"sync"

	"github.com/isayme/go-shadowsocks/shadowsocks/logger"
)

// Connection connection from client
type Connection struct {
	Remote *Remote
	Client *Client
}

// Serve start exchange data between remote & client
func (c Connection) Serve() {
	wg := sync.WaitGroup{}
	wg.Add(2)

	// any of remote/client closed, the other one should close with quiet
	closed := false

	go func() {
		defer wg.Done()

		_, err := copyBuffer(c.Remote, c.Client)
		if err != nil && !closed {
			logger.Errorf("io.Copy from client to remote fail, err: %#v", err)
		}
		closed = true
		logger.Debug("client read end")
		c.Remote.Close()
	}()

	go func() {
		defer wg.Done()

		_, err := copyBuffer(c.Client, c.Remote)
		if err != nil && !closed {
			logger.Errorf("io.Copy from remote to client fail, err: %#v", err)
		}
		closed = true
		logger.Debug("remote read end")
		c.Client.Close()
	}()

	wg.Wait()
}

var bufSize = os.Getpagesize()

// code from io.CopyBuffer, paste here to avoid buf escap to heap
func copyBuffer(dst io.Writer, src io.Reader) (written int64, err error) {
	buf := make([]byte, bufSize)

	for {
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

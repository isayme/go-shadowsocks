package connection

import (
	"io"
	"sync"

	"github.com/isayme/go-shadowsocks/shadowsocks/bufferpool"
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

		buf := bufferpool.DefaultPool.Get()
		defer bufferpool.DefaultPool.Put(buf)

		_, err := io.CopyBuffer(c.Remote, c.Client, buf)
		if err != nil && !closed {
			logger.Errorf("io.Copy from client to remote fail, err: %#v", err)
		}
		closed = true
		logger.Debug("client read end")
		c.Remote.Close()
	}()

	go func() {
		defer wg.Done()

		buf := bufferpool.DefaultPool.Get()
		defer bufferpool.DefaultPool.Put(buf)

		_, err := io.CopyBuffer(c.Client, c.Remote, buf)
		if err != nil && !closed {
			logger.Errorf("io.Copy from remote to client fail, err: %#v", err)
		}
		closed = true
		logger.Debug("remote read end")
		c.Client.Close()
	}()

	wg.Wait()
}

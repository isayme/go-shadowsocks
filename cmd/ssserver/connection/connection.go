package connection

import (
	"io"
	"sync"

	"github.com/isayme/go-shadowsocks/shadowsocks/logger"
)

// Connection connection from client
type Connection struct {
	Remote Remote
	Client *Client
}

// Serve start exchange data between remote & client
func (c Connection) Serve() {
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		// defer c.Remote.Close()
		_, err := io.Copy(c.Remote, c.Client)
		if err != nil {
			logger.Errorf("io.Copy from client to remote fail, err: %+v", err)
		}
	}()

	go func() {
		defer wg.Done()
		// defer c.Client.Close()
		_, err := io.Copy(c.Client, c.Remote)
		if err != nil {
			logger.Errorf("io.Copy from remote to client fail, err: %+v", err)
		}
	}()

	wg.Wait()
}

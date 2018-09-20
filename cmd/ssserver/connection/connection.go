package connection

import (
	"io"
	"os"
	"sync"

	"github.com/isayme/go-shadowsocks/shadowsocks/logger"
)

// BuffSize buf size for io read/write
var BuffSize = os.Getpagesize()

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
		buf := make([]byte, BuffSize)
		_, err := io.CopyBuffer(c.Remote, c.Client, buf)
		if err != nil {
			logger.Errorf("io.Copy from client to remote fail, err: %+v", err)
		}
	}()

	go func() {
		defer wg.Done()
		buf := make([]byte, BuffSize)
		_, err := io.CopyBuffer(c.Client, c.Remote, buf)
		if err != nil {
			logger.Errorf("io.Copy from remote to client fail, err: %+v", err)
		}
	}()

	wg.Wait()
}

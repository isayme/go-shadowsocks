package util

import (
	"io"
	"net"

	"github.com/panjf2000/ants"

	logger "github.com/isayme/go-logger"
	"github.com/pkg/errors"
)

func Proxy(client, server net.Conn) {
	defer client.Close()
	defer server.Close()

	// any of remote/client closed, the other one should close with quiet
	closed := false

	err := ants.Submit(func() {
		_, err := Copy(server, client)
		if err != nil && !closed {
			if errors.Cause(err) != io.EOF {
				logger.Errorf("[%s] Copy from client to server fail, err: %s", server.RemoteAddr(), err)
			}
		}
		closed = true
		server.Close()
		logger.Debug("client read end")
	})
	if err != nil {
		logger.Errorf("ants.Submit fail: %s", err)
		return
	}

	_, err = Copy(client, server)
	if err != nil && !closed {
		if errors.Cause(err) != io.EOF {
			logger.Errorf("[%s] Copy from server to client fail, err: %s", server.RemoteAddr(), err)
		}
	}
	closed = true
	client.Close()
	logger.Debug("remote read end")
}

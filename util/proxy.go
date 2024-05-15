package util

import (
	"net"
	"sync"

	"github.com/panjf2000/ants"

	logger "github.com/isayme/go-logger"
)

func Proxy(client net.Conn, tcpClinet *net.TCPConn, server net.Conn, tcpServer *net.TCPConn) {
	defer client.Close()
	defer server.Close()

	wg := sync.WaitGroup{}
	wg.Add(2)

	err := ants.Submit(func() {
		defer wg.Done()

		var err error
		_, err = Copy(server, client)
		tcpServer.CloseWrite()
		if err != nil {
			logger.Errorf("client read err: %v", err)
		} else {
			logger.Debug("client read end")
		}
	})
	if err != nil {
		logger.Errorf("ants.Submit fail: %s", err)
		return
	}

	err = ants.Submit(func() {
		defer wg.Done()

		var err error
		_, err = Copy(client, server)
		tcpClinet.CloseWrite()
		if err != nil {
			logger.Errorf("server read err: %v", err)
		} else {
			logger.Debug("server read end")
		}
	})
	if err != nil {
		logger.Errorf("ants.Submit fail: %s", err)
		return
	}

	wg.Wait()
	logger.Debug("proxy end")
}

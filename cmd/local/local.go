package local

import (
	"net"
	"strconv"
	"time"

	logger "github.com/isayme/go-logger"
	"github.com/isayme/go-shadowsocks/shadowsocks/conf"
	"github.com/isayme/go-shadowsocks/shadowsocks/ss"
	"github.com/panjf2000/ants"
	"github.com/pkg/errors"
)

func Run() {
	defer ants.Release()

	config := conf.Get()

	_ = logger.SetLevel(config.LogLevel)

	address := net.JoinHostPort(config.LocalAddress, strconv.Itoa(config.LocalPort))
	listener, err := net.Listen("tcp", address)
	if err != nil {
		logger.Panic(errors.Wrap(err, "net.Listen"))
	}
	logger.Infow("start listening", "address", address, "method", config.Method)

	timeout := time.Second * time.Duration(config.Timeout)
	client := ss.NewClient(config.Server, config.ServerPort, config.Method, config.Password, timeout)

	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Errorf("accept fail: %+v", err)
			continue
		}

		err = ants.Submit(func() error {
			client.AcceptAndHandle(conn)
			return nil
		})
		if err != nil {
			logger.Errorf("ants.Submit fail: %s", err)
			conn.Close()
		}
	}
}

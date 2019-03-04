package main

import (
	"flag"
	"net"
	"os"
	"strconv"
	"time"

	logger "github.com/isayme/go-logger"
	"github.com/isayme/go-shadowsocks/cmd/sslocal/socks5"
	"github.com/isayme/go-shadowsocks/shadowsocks/cipher"
	"github.com/isayme/go-shadowsocks/shadowsocks/conf"
	"github.com/isayme/go-shadowsocks/shadowsocks/util"
	"github.com/panjf2000/ants"
	"github.com/pkg/errors"
)

var showHelp = flag.Bool("h", false, "show help")
var showVersion = flag.Bool("v", false, "show version")

func main() {
	defer ants.Release()

	flag.Parse()

	if *showHelp {
		flag.Usage()
		os.Exit(0)
	}

	if *showVersion {
		util.PrintVersion()
		os.Exit(0)
	}

	config := conf.Get()

	_ = logger.SetLevel(config.LogLevel)

	address := net.JoinHostPort(config.LocalAddress, strconv.Itoa(config.LocalPort))
	listener, err := net.Listen("tcp", address)
	if err != nil {
		logger.Panic(errors.Wrap(err, "net.Listen"))
	}
	logger.Infow("start listening", "address", address, "method", config.Method)

	key := cipher.NewKey(config.Method, config.Password)
	timeout := time.Second * time.Duration(config.Timeout)

	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Errorf("accept fail: %+v", err)
			continue
		}

		cipher := cipher.NewCipher(config.Method)

		err = ants.Submit(func() error {
			handleConnection(conn, cipher, key, config.Server, config.ServerPort, timeout)
			return nil
		})
		if err != nil {
			logger.Errorf("ants.Submit fail: %s", err)
			conn.Close()
		}
	}
}

func handleConnection(conn net.Conn, cipher cipher.Cipher, key []byte, server string, serverPort int, timeout time.Duration) {
	client := util.NewTimeoutConn(conn, timeout)
	defer client.Close()

	logger.Debugf("new connection from: %s", client.RemoteAddr().String())

	request, err := socks5.NewRequest(client)
	if err != nil {
		logger.Errorf("NewRequest fail, err: %s", err)
		return
	}

	address := net.JoinHostPort(server, strconv.Itoa(int(serverPort)))
	ssconn, err := net.DialTimeout("tcp", address, 5*time.Second)
	if err != nil {
		logger.Errorf("dial ssserver fail, err: %s", err)
		return
	}
	cipher.Init(key, util.NewTimeoutConn(ssconn, timeout))
	defer cipher.Close()

	_, err = cipher.Write(request.RawAddr)
	if err != nil {
		logger.Errorf("write address fail, err: %s", err)
		return
	}

	util.Proxy(client, cipher)

	logger.Debugf("connection '%s' closed", request.RemoteAddress())
}

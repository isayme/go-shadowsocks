package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/isayme/go-logger"
	"github.com/isayme/go-shadowsocks/cmd/ssserver/connection"
	"github.com/isayme/go-shadowsocks/shadowsocks/aead"
	"github.com/isayme/go-shadowsocks/shadowsocks/cipher"
	"github.com/isayme/go-shadowsocks/shadowsocks/conf"
	"github.com/panjf2000/ants"
	"github.com/pkg/errors"
)

var configPath = flag.String("c", "/etc/shadowsocks.json", "config file path")
var showHelp = flag.Bool("h", false, "show help")
var showVersion = flag.Bool("v", false, "show version")

// Version current version
var Version = "unkonwn"

func main() {
	defer ants.Release()

	flag.Parse()

	if *showHelp {
		flag.Usage()
		os.Exit(0)
	}

	if *showVersion {
		fmt.Printf("%s: %s\n", os.Args[0], Version)
		os.Exit(0)
	}

	config := conf.Get()

	logger.SetLevel(config.LogLevel)

	address := fmt.Sprintf("%s:%d", config.Server, config.ServerPort)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		logger.Panic(errors.Wrap(err, "net.Listen"))
	}

	logger.Infow("start listening", "address", address, "method", config.Method)

	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Errorf("accept fail: %+v", err)
			continue
		}

		c, err := aead.NewCipher(config.Method, config.Password, conn)
		if err != nil {
			logger.Panic(errors.Wrap(err, "create cipher"))
		}

		ants.Submit(func() error {
			handleConnection(conn, c, config.Timeout)
			return nil
		})
	}
}

func handleConnection(conn net.Conn, c cipher.Cipher, timeout int) {
	logger.Debugf("new connection from: %s", conn.RemoteAddr().String())

	client, err := connection.NewClient(conn, c)
	if err != nil {
		logger.Errorf("NewClient fail, err: %+v", err)
		return
	}
	defer client.Close()

	// read address type
	address, err := client.ReadAddress(timeout)
	if err != nil {
		logger.Errorf("read address fail, err: %+v", err)
		return
	}

	logger.Infof("connecting remote [%s]", address)
	remote, err := net.Dial("tcp", address)
	if err != nil {
		logger.Warnf("dial remote [%s] failed, err: %+v", address, err)
		return
	}
	defer remote.Close()

	logger.Debugf("connect remote [%s] success", address)

	connection := connection.Connection{
		Client:  client,
		Remote:  connection.NewRemote(remote),
		Timeout: timeout,
	}

	connection.Serve()

	logger.Debugf("connection [%s] closed", address)
}

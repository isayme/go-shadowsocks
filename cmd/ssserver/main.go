package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/isayme/go-shadowsocks/cmd/ssserver/connection"
	"github.com/isayme/go-shadowsocks/shadowsocks/cipher"
	"github.com/isayme/go-shadowsocks/shadowsocks/conf"
	"github.com/isayme/go-shadowsocks/shadowsocks/logger"
	"github.com/pkg/errors"
)

var configPath = flag.String("c", "/etc/shadowsocks.json", "config file path")
var showHelp = flag.Bool("h", false, "show help")
var showVersion = flag.Bool("v", false, "show version")

// Version current version
var Version = "unkonwn"

func main() {
	flag.Parse()

	if *showHelp {
		flag.Usage()
		os.Exit(0)
	}

	if *showVersion {
		fmt.Printf("%s: %s\n", os.Args[0], Version)
		os.Exit(0)
	}

	config, err := conf.ParseConfig(*configPath)
	if err != nil {
		logger.Panic(errors.Wrap(err, "parse config"))
	}

	address := fmt.Sprintf("%s:%d", config.Server, config.ServerPort)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		logger.Panic(errors.Wrap(err, "net.Listen"))
	}

	logger.Infof("start listening %s", address)

	c, err := cipher.NewCipher(config.Method, config.Password)
	if err != nil {
		logger.Panic(errors.Wrap(err, "create cipher"))
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Errorf("accept fail: %+v", err)
			continue
		}

		go handleConnection(conn, *c)
	}
}

func handleConnection(conn net.Conn, c cipher.Cipher) {
	logger.Debugf("new connection from: %s", conn.RemoteAddr().String())

	client, err := connection.NewClient(conn, c)
	if err != nil {
		logger.Errorf("NewClient fail, err: %+v", err)
		return
	}

	// read address type
	address, err := client.ReadAddress()
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
	logger.Debugf("connect remote [%s] success", address)

	connection := connection.Connection{
		Client: client,
		Remote: remote,
	}

	connection.Serve()

	logger.Debugf("connection [%s] closed", address)
}

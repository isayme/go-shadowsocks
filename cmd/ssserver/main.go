package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"github.com/isayme/go-logger"
	"github.com/isayme/go-shadowsocks/cmd/ssserver/connection"
	"github.com/isayme/go-shadowsocks/shadowsocks/bufferpool"
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

	logger.SetLevel(config.LogLevel)

	address := fmt.Sprintf("%s:%d", config.Server, config.ServerPort)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		logger.Panic(errors.Wrap(err, "net.Listen"))
	}

	logger.Infow("start listening", "address", address, "method", config.Method)

	key := cipher.NewKey(config.Method, config.Password)

	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Errorf("accept fail: %+v", err)
			continue
		}

		c := cipher.NewCipher(config.Method)
		c.Init(key, conn)

		ants.Submit(func() error {
			handleConnection(conn, c, config.Timeout)
			return nil
		})
	}
}

const handleshakeTimeout = 5 // second

func handleConnection(conn net.Conn, c cipher.Cipher, timeout int) {
	logger.Debugf("new connection from: %s", conn.RemoteAddr().String())

	client, err := connection.NewClient(conn, c)
	if err != nil {
		logger.Errorf("NewClient fail, err: %+v", err)
		return
	}
	defer client.Close()

	// read address type
	address, err := client.ReadAddress(handleshakeTimeout) // 5s timeout for address
	if err != nil {
		logger.Errorw("read address fail", "err", err, "remoteAddr", conn.RemoteAddr().String())

		// random response
		buf := bufferpool.Get(16)
		defer bufferpool.Put(buf)
		io.ReadFull(rand.Reader, buf)
		client.Write(buf)

		return
	}

	logger.Infof("connecting remote [%s]", address)
	// dial with timeout
	remote, err := net.DialTimeout("tcp", address, handleshakeTimeout*time.Second)
	if err != nil {
		logger.Warnf("dial remote [%s] failed, err: %+v", address, err)
		return
	}
	defer remote.Close()

	logger.Debugf("connect remote [%s] success", address)

	connection := connection.Connection{
		Client:  client,
		Remote:  connection.NewRemote(remote, address),
		Timeout: timeout,
	}

	connection.Serve()

	logger.Debugf("connection [%s] closed", address)
}

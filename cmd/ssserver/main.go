package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"

	"encoding/binary"
	"io"

	logger "github.com/isayme/go-logger"
	"github.com/isayme/go-shadowsocks/shadowsocks/cipher"
	"github.com/isayme/go-shadowsocks/shadowsocks/conf"
	"github.com/isayme/go-shadowsocks/shadowsocks/util"
	"github.com/panjf2000/ants"
	"github.com/pkg/errors"

	"github.com/isayme/go-shadowsocks/shadowsocks/bufferpool"
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

	address := fmt.Sprintf("%s:%d", config.Server, config.ServerPort)
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
			handleConnection(conn, cipher, key, timeout)
			return nil
		})
		if err != nil {
			logger.Errorf("ants.Submit fail: %s", err)
			conn.Close()
		}
	}
}

func handleConnection(conn net.Conn, cipher cipher.Cipher, key []byte, timeout time.Duration) {
	conn = util.NewTimeoutConn(conn, timeout)
	defer cipher.Close()

	logger.Debugf("new connection from: %s", conn.RemoteAddr().String())
	cipher.Init(key, conn)

	// read address type
	address, err := readAddress(cipher)
	if err != nil {
		logger.Errorw("read address fail", "err", err, "remoteAddr", conn.RemoteAddr().String())
		return
	}

	logger.Infof("connecting remote [%s]", address)
	// dial with timeout
	remote, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		logger.Warnf("dial remote [%s] failed, err: %+v", address, err)
		return
	}
	remote = util.NewTimeoutConn(remote, timeout)
	defer remote.Close()

	logger.Debugf("connect remote [%s] success", address)

	util.Proxy(cipher, remote)

	logger.Debugf("connection [%s] closed", address)
}

func readAddress(r io.Reader) (string, error) {
	data := bufferpool.Get(256)
	defer bufferpool.Put(data)

	if _, err := io.ReadFull(r, data[:1]); err != nil {
		return "", errors.Wrap(err, "read type")
	}

	typ := data[0]
	logger.Debugf("address type: %02x", typ)

	var host string
	switch typ {
	case util.AddressTypeIPV4:
		if _, err := io.ReadFull(r, data[:net.IPv4len]); err != nil {
			return "", errors.Wrap(err, "read ipv4")
		}
		host = net.IP(data[:net.IPv4len]).String()
	case util.AddressTypeDomain:
		if _, err := io.ReadFull(r, data[:1]); err != nil {
			return "", errors.Wrap(err, "read domain length")
		}
		domainLen := int(data[0])

		if _, err := io.ReadFull(r, data[:domainLen]); err != nil {
			return "", errors.Wrap(err, "read domain")
		}
		host = string(data[:domainLen])
	case util.AddressTypeIPV6:
		if _, err := io.ReadFull(r, data[:net.IPv6len]); err != nil {
			return "", errors.Wrap(err, "read ipv6")
		}
		host = net.IP(data[:net.IPv6len]).String()
	default:
		return "", errors.Errorf("invalid address type: %02x", typ)
	}

	if _, err := io.ReadFull(r, data[:2]); err != nil {
		return "", errors.Wrap(err, "read port")
	}

	port := binary.BigEndian.Uint16(data)

	return fmt.Sprintf("%s:%d", host, port), nil
}

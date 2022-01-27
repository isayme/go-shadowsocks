package aead

import (
	"bytes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"net"

	"github.com/isayme/go-bufferpool"
	"github.com/pkg/errors"
)

// cipherInfo cipher definition
type cipherInfo struct {
	keySize int

	genEncryptAEAD func(key, salt []byte, keySize int) (cipher.AEAD, error)
	genDecryptAEAD func(key, salt []byte, keySize int) (cipher.AEAD, error)
}

func (ci *cipherInfo) getSaltSize() int {
	return ci.keySize
}

// Cipher cipher
type Cipher struct {
	enc cipher.AEAD
	dec cipher.AEAD

	key []byte

	buffer *bytes.Buffer

	rnonce []byte
	wnonce []byte

	info *cipherInfo
}

// NewCipher create aead cipher
func NewCipher(method string, key []byte) *Cipher {
	c := &Cipher{}
	c.key = key

	info, ok := cipherMethods[method]
	if !ok {
		panic(fmt.Errorf("unsupported method: %s", method))
	}

	c.info = info

	c.buffer = bytes.NewBuffer(nil)

	return c
}

// keySize return key size
func (c *Cipher) KeySize() int {
	return c.info.keySize
}

func (c *Cipher) getEncryptAEAD(salt []byte) (s cipher.AEAD, err error) {
	_, err = io.ReadFull(rand.Reader, salt)
	if err != nil {
		return nil, err
	}

	s, err = c.info.genEncryptAEAD(c.key, salt, c.info.keySize)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (c *Cipher) getDecryptAEAD(salt []byte) (cipher.AEAD, error) {
	return c.info.genDecryptAEAD(c.key, salt, c.info.keySize)
}

func (c *Cipher) read(conn net.Conn) error {
	// read size
	sizeBuf := bufferpool.Get(c.dec.Overhead() + 2)
	defer bufferpool.Put(sizeBuf)

	_, err := io.ReadFull(conn, sizeBuf)
	if err != nil {
		return errors.Wrap(err, "aead read size")
	}

	ret, err := c.dec.Open(sizeBuf[:0], c.rnonce, sizeBuf, nil)
	if err != nil {
		return errors.Wrap(err, "aead decrypt size")
	}
	increment(c.rnonce)
	payloadSize := int(binary.BigEndian.Uint16(ret)&0x3FFF) + c.dec.Overhead()

	// read payload
	payloadBuf := bufferpool.Get(payloadSize)
	defer bufferpool.Put(payloadBuf)

	_, err = io.ReadFull(conn, payloadBuf)
	if err != nil {
		return errors.Wrap(err, "aead read payload")
	}

	ret, err = c.dec.Open(payloadBuf[:0], c.rnonce, payloadBuf, nil)
	if err != nil {
		return errors.Wrap(err, "aead decrypt payload")
	}
	increment(c.rnonce)
	c.buffer.Write(ret)

	return nil
}

// Read read from client
func (c *Cipher) Read(conn net.Conn, p []byte) (n int, err error) {
	if c.dec == nil {
		salt := bufferpool.Get(c.info.getSaltSize())
		defer bufferpool.Put(salt)

		if _, err = io.ReadFull(conn, salt); err != nil {
			return 0, err
		}

		s, err := c.getDecryptAEAD(salt)
		if err != nil {
			return 0, err
		}

		c.dec = s

		// init read nonce
		c.rnonce = make([]byte, s.NonceSize())
	}

	if c.buffer.Len() > 0 {
		return c.buffer.Read(p)
	}

	err = c.read(conn)
	if err != nil {
		return 0, err
	}

	return c.buffer.Read(p)
}

func (c *Cipher) encrypt(dst, src []byte) (n int) {
	size := len(src)

	binary.BigEndian.PutUint16(dst, uint16(size))

	ret := c.enc.Seal(dst[:0], c.wnonce, dst[:2], nil)
	increment(c.wnonce)
	n = len(ret)

	ret = c.enc.Seal(dst[n:n], c.wnonce, src, nil)
	increment(c.wnonce)
	n += len(ret)

	return n
}

// Write write to client
func (c *Cipher) Write(conn net.Conn, p []byte) (n int, err error) {
	if c.enc == nil {
		salt := bufferpool.Get(c.info.getSaltSize())
		defer bufferpool.Put(salt)

		c.enc, err = c.getEncryptAEAD(salt)
		if err != nil {
			return 0, err
		}

		nw, err := conn.Write(salt)
		if err != nil {
			return 0, errors.Wrap(err, "salt write")
		}
		if nw != len(salt) {
			return 0, errors.New("salt short write")
		}

		// init write nonce
		c.wnonce = make([]byte, c.enc.NonceSize())
	}

	size := len(p)

	buf := bufferpool.Get(c.enc.Overhead() + 2 + size + c.enc.Overhead())
	defer bufferpool.Put(buf)

	n = c.encrypt(buf, p)

	nw, err := conn.Write(buf[:n])
	if err != nil {
		return 0, errors.Wrapf(err, "write short: %d/%d", nw, n)
	}

	return size, nil
}

func increment(b []byte) {
	for i := range b {
		b[i]++
		if b[i] != 0 {
			return
		}
	}
}

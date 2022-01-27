package aead

import (
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/isayme/go-bufferpool"
	"github.com/isayme/go-shadowsocks/util"
	"github.com/pkg/errors"
)

type aeadWriter struct {
	key []byte

	writer    io.Writer
	newCipher func([]byte) (cipher.AEAD, error)

	aead  cipher.AEAD
	nonce []byte
}

func NewWriter(writer io.Writer, key []byte, newCipher func([]byte) (cipher.AEAD, error)) *aeadWriter {
	return &aeadWriter{
		key:       key,
		writer:    writer,
		newCipher: newCipher,
	}
}

func (w *aeadWriter) getAeadCipher() (cipher.AEAD, error) {
	if w.aead != nil {
		return w.aead, nil
	}

	salt := bufferpool.Get(len(w.key))
	defer bufferpool.Put(salt)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, err
	}

	key := bufferpool.Get(len(w.key))
	defer bufferpool.Put(key)
	hkdfSHA1(w.key, salt, hkdfInfo, key)
	c, err := w.newCipher(key)
	if err != nil {
		return nil, errors.Wrap(err, "new cipher")
	}

	w.aead = c
	w.nonce = make([]byte, c.NonceSize())

	n, err := w.writer.Write(salt)
	if err != nil {
		return nil, err
	}
	if n != len(salt) {
		return nil, fmt.Errorf("write salt short")
	}

	return w.aead, nil
}

func (w *aeadWriter) Write(p []byte) (n int, err error) {
	c, err := w.getAeadCipher()
	if err != nil {
		return 0, err
	}

	size := len(p)

	buf := bufferpool.Get(2 + c.Overhead() + size + c.Overhead())
	defer bufferpool.Put(buf)

	// write size
	binary.BigEndian.PutUint16(buf, uint16(size))
	ret := c.Seal(buf[:0], w.nonce, buf[:2], nil)
	util.NextNonce(w.nonce)
	n = len(ret)

	// write payload
	ret = c.Seal(buf[n:n], w.nonce, p, nil)
	util.NextNonce(w.nonce)
	n += len(ret)

	nw, err := w.writer.Write(buf[:n])
	if err != nil {
		return 0, err
	}

	if nw != n {
		return 0, errors.New("write short")
	}

	return size, nil
}

package aead

import (
	"bytes"
	"crypto/cipher"
	"encoding/binary"
	"io"

	"github.com/isayme/go-bufferpool"
	"github.com/isayme/go-shadowsocks/util"
	"github.com/pkg/errors"
)

type aeadReader struct {
	key []byte

	reader     io.Reader
	newCipher  func([]byte) (cipher.AEAD, error)
	aead       cipher.AEAD
	nonce      []byte
	readBuffer *bytes.Buffer
}

func NewReader(reader io.Reader, key []byte, newCipher func([]byte) (cipher.AEAD, error)) *aeadReader {
	return &aeadReader{
		key:        key,
		reader:     reader,
		readBuffer: bytes.NewBuffer(nil),
		newCipher:  newCipher,
	}
}

func (r *aeadReader) getAeadCipher() (cipher.AEAD, error) {
	if r.aead != nil {
		return r.aead, nil
	}

	salt := bufferpool.Get(len(r.key))
	defer bufferpool.Put(salt)

	if _, err := io.ReadFull(r.reader, salt); err != nil {
		return nil, errors.Wrap(err, "read aead salt")
	}

	key := bufferpool.Get(len(r.key))
	defer bufferpool.Put(key)
	hkdfSHA1(r.key, salt, hkdfInfo, key)
	c, err := r.newCipher(key)
	if err != nil {
		return nil, errors.Wrap(err, "new cipher")
	}

	r.aead = c
	r.nonce = make([]byte, c.NonceSize())

	return r.aead, nil
}

func (r *aeadReader) doRead() error {
	c, err := r.getAeadCipher()
	if err != nil {
		return err
	}

	// read size
	sizeBuf := bufferpool.Get(2 + c.Overhead())
	defer bufferpool.Put(sizeBuf)
	_, err = io.ReadFull(r.reader, sizeBuf)
	if err != nil {
		return errors.Wrap(err, "aead read size")
	}

	ret, err := c.Open(sizeBuf[:0], r.nonce, sizeBuf, nil)
	if err != nil {
		return errors.Wrap(err, "aead decrypt size")
	}
	util.NextNonce(r.nonce)

	payloadSize := int(binary.BigEndian.Uint16(ret))

	// read payload
	payloadBuf := bufferpool.Get(payloadSize + c.Overhead())
	defer bufferpool.Put(payloadBuf)

	_, err = io.ReadFull(r.reader, payloadBuf)
	if err != nil {
		return errors.Wrap(err, "aead read payload")
	}

	ret, err = c.Open(payloadBuf[:0], r.nonce, payloadBuf, nil)
	if err != nil {
		return errors.Wrap(err, "aead decrypt payload")
	}
	util.NextNonce(r.nonce)

	r.readBuffer.Write(ret)

	return nil
}

func (r *aeadReader) Read(p []byte) (n int, err error) {
	if r.readBuffer.Len() > 0 {
		return r.readBuffer.Read(p)
	}

	if err := r.doRead(); err != nil {
		return 0, err
	}

	return r.readBuffer.Read(p)
}

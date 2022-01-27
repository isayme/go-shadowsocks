package stream

import (
	"crypto/cipher"
	"crypto/rand"
	"io"

	"github.com/isayme/go-bufferpool"
)

type streamWriter struct {
	key    []byte
	ivSize int

	writer io.Writer

	newCipher func([]byte, []byte) (cipher.Stream, error)

	stream cipher.Stream
}

func NewWriter(writer io.Writer, key []byte, ivSize int, newCipher func([]byte, []byte) (cipher.Stream, error)) *streamWriter {
	return &streamWriter{
		key:       key,
		ivSize:    ivSize,
		writer:    writer,
		newCipher: newCipher,
	}
}

func (r *streamWriter) Write(p []byte) (n int, err error) {
	if r.stream == nil {
		iv := bufferpool.Get(r.ivSize)
		defer bufferpool.Put(iv)
		if _, err = io.ReadFull(rand.Reader, iv); err != nil {
			return 0, err
		}

		s, err := r.newCipher(r.key, iv)
		if err != nil {
			return 0, err
		}

		r.writer.Write(iv)
		r.stream = s
	}

	r.stream.XORKeyStream(p, p)
	n, err = r.writer.Write(p)
	return n, err
}

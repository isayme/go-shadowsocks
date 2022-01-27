package stream

import (
	"crypto/cipher"
	"io"

	"github.com/isayme/go-bufferpool"
)

type streamReader struct {
	key    []byte
	ivSize int

	reader io.Reader

	newCipher func([]byte, []byte) (cipher.Stream, error)

	stream cipher.Stream
}

func NewReader(reader io.Reader, key []byte, ivSize int, newCipher func([]byte, []byte) (cipher.Stream, error)) *streamReader {
	return &streamReader{
		key:       key,
		ivSize:    ivSize,
		reader:    reader,
		newCipher: newCipher,
	}
}

func (r *streamReader) Read(p []byte) (n int, err error) {
	if r.stream == nil {
		iv := bufferpool.Get(r.ivSize)
		defer bufferpool.Put(iv)
		if _, err = io.ReadFull(r.reader, iv); err != nil {
			return 0, err
		}

		s, err := r.newCipher(r.key, iv)
		if err != nil {
			return 0, err
		}

		r.stream = s
	}

	n, err = r.reader.Read(p)
	r.stream.XORKeyStream(p, p[0:n])
	return n, err
}

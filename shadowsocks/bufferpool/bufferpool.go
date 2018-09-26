package bufferpool

import (
	"os"
	"sync"
)

// Pool pool
type Pool struct {
	size int
	pool *sync.Pool
}

// DefaultPool default pool
var DefaultPool *Pool

// NewPool create pool
func NewPool(size int) *Pool {
	return &Pool{
		size: size,
		pool: &sync.Pool{
			New: func() interface{} {
				return make([]byte, size)
			},
		},
	}
}

// Get get a buf
func (p *Pool) Get() []byte {
	return p.pool.Get().([]byte)
}

// Put release a buf
func (p *Pool) Put(x []byte) {
	p.pool.Put(x)
}

func init() {
	size := os.Getpagesize()
	DefaultPool = NewPool(size)
}

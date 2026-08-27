package chunk

import "sync"

type Pool struct {
	size int
	pool sync.Pool
}

func NewPool(size int) *Pool {
	return &Pool{size: size, pool: sync.Pool{New: func() any { return make([]byte, size) }}}
}

func (p *Pool) Get() []byte { return p.pool.Get().([]byte)[:p.size] }
func (p *Pool) Put(b []byte) {
	if cap(b) >= p.size {
		p.pool.Put(b[:p.size])
	}
}

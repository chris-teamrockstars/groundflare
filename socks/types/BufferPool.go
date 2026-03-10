package types

import "sync"

type BufferPool struct {
	size int
	pool *sync.Pool
}

func NewBufferPool(size int) *BufferPool {

	return &BufferPool{
		size: size,
		pool: &sync.Pool{
			New: func() interface{} {
				return make([]byte, 0, size)
			},
		},
	}

}

func (buffer_pool *BufferPool) Get() []byte {
	return buffer_pool.pool.Get().([]byte)
}

func (buffer_pool *BufferPool) Put(data []byte) {

	if cap(data) != buffer_pool.size {
		panic("invalid buffer size that's put into leaky buffer")
	}

	// This ensures cap=size, length=0
	buffer_pool.pool.Put(data[:0])

}

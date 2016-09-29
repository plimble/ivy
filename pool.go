package ivy

import (
	"bytes"
)

type bufferPool struct {
	list chan *bytes.Buffer
}

func newBufferPool(poolSize int) *bufferPool {
	b := &bufferPool{
		list: make(chan *bytes.Buffer, poolSize),
	}

	return b
}

func (p *bufferPool) Put(b *bytes.Buffer) {
	b.Reset()
	select {
	case p.list <- b:
	default:
	}
}

func (p *bufferPool) Get() *bytes.Buffer {
	select {
	case b := <-p.list:
		return b
	default:
		return &bytes.Buffer{}
	}
}

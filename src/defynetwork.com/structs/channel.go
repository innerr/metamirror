package structs

import "io"

func NewBiChannel() *BiChannel {
	p := &BiChannel{&MemChannel{}, &MemChannel{}}
	p.A.fsend = func(r io.Reader, n uint32) {
		p.B.freceive(r, n)
	}
	p.B.fsend = func(r io.Reader, n uint32) {
		p.A.freceive(r, n)
	}
	return p
}

type BiChannel struct {
	A *MemChannel
	B *MemChannel
}

func (p *MemChannel) Send(r io.Reader, n uint32) {
	p.fsend(r, n)
}

func (p *MemChannel) Receive(fun func(io.Reader, uint32)) {
	p.freceive = fun
}

type MemChannel struct {
	fsend func(io.Reader, uint32)
	freceive func(io.Reader, uint32)
}

type IChannel interface {
	Send(io.Reader, uint32)
	Receive(func(io.Reader, uint32))
}

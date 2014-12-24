package structs

import "io"

func (p *BiChannel) Close() {
	p.A.Close()
	p.B.Close()
}

func NewBiChannel() *BiChannel {
	a := NewMemChannel(nil, nil)
	b := NewMemChannel(nil, nil)
	a.fsend = func(r io.Reader, n uint32) {
		b.freceive(r, n)
	}
	b.fsend = func(r io.Reader, n uint32) {
		a.freceive(r, n)
	}
	return &BiChannel{NewAsynChannel(a, 0), NewAsynChannel(b, 0)}
}

type BiChannel struct {
	A IChannel
	B IChannel
}

func (p *MemChannel) Close() {
}

func (p *MemChannel) Send(r io.Reader, n uint32) {
	p.fsend(r, n)
}

func (p *MemChannel) Receive(fun Transport) {
	p.freceive = fun
	if p.fsend == nil {
		p.fsend = fun
	}
}

func NewMemChannel(send, recv Transport) *MemChannel {
	if send == nil {
		send = recv
	}
	return &MemChannel{send, recv}
}

type MemChannel struct {
	fsend Transport
	freceive Transport
}

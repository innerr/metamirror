package network

import (
	"bytes"
	"encoding/binary"
	"io"
	"net"
	"sync"
	"defynetwork.com/tools"
)

func (p *TcpChannel) Start()  {
	for {
		n := uint32(0)
		err := binary.Read(p.conn, binary.LittleEndian, &n)
		if p.closed(err) {
			break
		}
		if p.frecv == nil {
			panic("no recv func")
		}

p.preread = true
		if !p.preread {
			p.frecv(p.conn, n)
			continue
		}

		buf := new(bytes.Buffer)
		_, err = io.CopyN(buf, p.conn, int64(n))
		if p.closed(err) {
			break
		}
		data := buf.Bytes()
		p.frecv(bytes.NewReader(data), n)
	}
}

func (p *TcpChannel) Send(r io.Reader, n uint32) {
	p.lock.Lock()
	defer p.lock.Unlock()
	func () {
		defer func() {
			e := recover()
			if e != nil {
				err, ok := e.(error)
				if !ok {
					panic(err)
				}
				p.closed(err)
			}
		}()
		tools.Dump(p.conn, n)
	}()
	_, err := io.CopyN(p.conn, r, int64(n))
	p.closed(err)
	return
}

func (p *TcpChannel) closed(err error) bool {
	if err == nil {
		return false
	}
	if tools.NetworkErr(err) || err.Error() == "EOF" {
		p.frecv = nil
		for _, it := range p.fcs {
			it(p.conn)
		}
		return true
	}
	panic(err)
}

func (p *TcpChannel) Receive(fun func(io.Reader, uint32))  {
	if p.frecv != nil && fun != nil {
		panic("double assigned")
	}
	p.frecv = fun
}

func (p *TcpChannel) OnClose(fun func(net.Conn)) {
	if fun == nil {
		panic("nil func")
	}
	p.fcs = append(p.fcs, fun)
}

func NewTcpChannel(log *tools.Log, conn net.Conn, preread bool) *TcpChannel {
	return &TcpChannel{log, conn, nil, nil, preread, sync.Mutex{}}
}

type TcpChannel struct {
	log *tools.Log
	conn net.Conn
	frecv func(io.Reader, uint32)
	fcs []func(net.Conn)
	preread bool
	lock sync.Mutex
}

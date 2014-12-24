package structs

import (
	"bytes"
	"io"
	"strconv"
	"defynetwork.com/tools"
)

func (p *Rpc) Func(name string) *RpcFunction{
	fun, ok := p.funs[name]
	if !ok {
		fun = &RpcFunction{p, name}
		p.funs[name] = fun
	}
	return fun
}

func (p *Rpc) register(name string, fun Transport) {
	_, ok := p.frecvs[name]
	if ok {
		panic("double assigned")
	}
	p.frecvs[name] = fun
}

func (p *Rpc) send(name string, r io.Reader, n uint32) {
	if len(name) == 0 {
		panic("no sub name")
	}
	buf := new(bytes.Buffer)
	tools.Dumps(buf, name)
	data := buf.Bytes()
	p.ch.Send(io.MultiReader(bytes.NewReader(data), r), n + uint32(len(data)))
}

func (p *Rpc) receive(r io.Reader, n uint32) {
	name := tools.Loads(r)
	fun, ok := p.frecvs[name]
	if !ok {
		err := name
		if len(name) > 64 {
			err = name[:64] + "..." + strconv.Itoa(len(name))
		}
		panic("handler not exists: " + err)
	}
	fun(r, n - uint32(tools.DumpsSize(name)))
}

func (p *Rpc) Close() {
	p.ch.Close()
	for name, _ := range p.frecvs {
		delete(p.frecvs, name)
	}
}

func NewRpc(ch IChannel) *Rpc {
	p := &Rpc{ch, make(map[string]Transport), make(map[string]*RpcFunction)}
	p.ch.Receive(p.receive)
	return p
}

type Rpc struct {
	ch IChannel
	frecvs map[string]Transport
	funs map[string]*RpcFunction
}

func (p *RpcFunction) Receive(fun Transport) {
	p.rpc.register(p.name, fun)
}

func (p *RpcFunction) Send(r io.Reader, n uint32) {
	p.rpc.send(p.name, r, n)
}

type RpcFunction struct {
	rpc *Rpc
	name string
}

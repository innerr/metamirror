package structs

import (
	"bytes"
	"io"
	"strconv"
	"defynetwork.com/tools"
)

func (p *Channels) Sub(name string) *_SubChannel {
	child, ok := p.children[name]
	if !ok {
		child = &_SubChannel{p, name}
		p.children[name] = child
	}
	return child
}

func (p *Channels) Send(name string, r io.Reader, n uint32) {
	if len(name) == 0 {
		panic("no sub name")
	}
	buf := new(bytes.Buffer)
	tools.Dumps(buf, name)
	data := buf.Bytes()
	p.origin.Send(io.MultiReader(bytes.NewReader(data), r), n + uint32(len(data)))
}

func (p *Channels) Clean(name string) {
	delete(p.frecvs, name)
}

func (p *Channels) CleanAll() {
	for name, _ := range p.frecvs {
		delete(p.frecvs, name)
	}
}

func (p *Channels) Receive(name string, fun func(io.Reader, uint32)) {
	_, ok := p.frecvs[name]
	if ok {
		panic("double assigned")
	}
	p.frecvs[name] = fun
}

func (p *Channels) receive(r io.Reader, n uint32) {
	name := tools.Loads(r)
	recv, ok := p.frecvs[name]
	if !ok {
		err := name
		if len(name) > 64 {
			err = name[:64] + "..." + strconv.Itoa(len(name))
		}
		panic("handler not exists: " + err)
	}
	recv(r, n - 2 - uint32(len(name)))
}

func NewChannels(origin IChannel) *Channels {
	p := &Channels{origin, make(_RecFuncs), make(map[string]*_SubChannel)}
	p.origin.Receive(p.receive)
	return p
}

type Channels struct {
	origin IChannel
	frecvs _RecFuncs
	children map[string]*_SubChannel
}

type _RecFuncs map[string]func(io.Reader, uint32)

func (p *_SubChannel) Receive(fun func(io.Reader, uint32)) {
	if fun == nil {
		p.parent.Clean(p.name)
	} else {
		p.parent.Receive(p.name, fun)
	}
}

func (p *_SubChannel) Send(r io.Reader, n uint32) {
	p.parent.Send(p.name, r, n)
}

type _SubChannel struct {
	parent *Channels
	name string
}

package network

import (
	"io"
	"bytes"
	"defynetwork.com/tools"
)

func (p *CodedChannel) DoSend() {
	for _, data := range p.cache {
		p.ch.Send(bytes.NewReader(data), uint32(len(data)))
	}
	p.cache = nil
}

func (p *CodedChannel) Send(r io.Reader, n uint32) {
	data := tools.Encode(r, n, p.compress, p.keyword)
	if !p.lazy {
		p.ch.Send(bytes.NewReader(data), uint32(len(data)))
	} else {
		p.cache = append(p.cache, data)
	}
}

func (p *CodedChannel) Receive(fun func(io.Reader, uint32)) {
	p.ch.Receive(func(r io.Reader, n uint32) {
		data := tools.ReadN(r, n)
		data = tools.Decode(data, p.keyword)
		fun(bytes.NewReader(data), uint32(len(data)))
	})
}

func NewCodedChannel(ch IChannel, compress bool, lazy bool, keyword []byte) *CodedChannel {
	return &CodedChannel{ch, compress, lazy, keyword, nil}
}

type CodedChannel struct {
	ch IChannel
	compress bool
	lazy bool
	keyword []byte
	cache [][]byte
}

func NormalizeKeyword(keyword string) []byte {
	data := []byte(keyword)
	size := len(data)
	if size > 32 {
		data = data[:32]
	}
	if size < 32 {
		data = append(data, make([]byte, 32 - size)...)
	}
	return data
}

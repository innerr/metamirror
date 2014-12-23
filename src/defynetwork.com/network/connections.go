package network

import (
	"bytes"
	"io"
	"net"
	"sync"
	"defynetwork.com/tools"
)

func (p *Connections) Walk(fun func(*ConnInfo)) {
	p.lock.Lock()
	defer p.lock.Unlock()
	for _, info := range p.mconns {
		fun(info)
	}
}

func (p *Connections) ByHost(host string) *ConnInfo {
	p.lock.Lock()
	defer p.lock.Unlock()
	info, ok := p.hosts[host]
	if !ok {
		return nil
	}
	return info
}

func (p *Connections) TransConns() []*Channels {
	p.lock.Lock()
	defer p.lock.Unlock()
	chss := []*Channels{}
	for chs, _ := range p.tconns {
		chss = append(chss, chs)
	}
	return chss
}

func (p *Connections) DelTransConn(chs *Channels) {
	p.lock.Lock()
	defer p.lock.Unlock()
	delete(p.tconns, chs)
}

func (p *Connections) AddTransConn(conn net.Conn, chs *Channels, host string) {
	p.lock.Lock()
	defer p.lock.Unlock()
	info, ok := p.hosts[host]
	if !ok {
		panic("main conn losted")
	}
	info.TransConn = conn
	info.TransChs = chs
	p.tconns[chs] = info
}

func (p *Connections) ByTrans(chs *Channels) *ConnInfo {
	p.lock.Lock()
	defer p.lock.Unlock()
	info, ok := p.tconns[chs]
	if !ok {
		return nil
	}
	return info
}

func (p *Connections) ByHid(hid uint64) *ConnInfo {
	p.lock.Lock()
	defer p.lock.Unlock()
	info, ok := p.hids[hid]
	if !ok {
		return nil
	}
	return info
}

func (p *Connections) ByMain(chs *Channels) *ConnInfo {
	p.lock.Lock()
	defer p.lock.Unlock()
	info, ok := p.mconns[chs]
	if !ok {
		return nil
	}
	return info
}

func (p *Connections) Close(info *ConnInfo) {
	p.lock.Lock()
	defer p.lock.Unlock()
	info.MainConn.Close()
	if info.TransConn != nil {
		info.TransConn.Close()
	}
}

func (p *Connections) Size() int {
	p.lock.Lock()
	defer p.lock.Unlock()
	return len(p.mconns)
}

func (p *Connections) Unbind(conn net.Conn, chs *Channels) {
	p.lock.Lock()
	defer p.lock.Unlock()
	info, ok := p.mconns[chs]
	if !ok {
		return
	}
	delete(p.mconns, chs)
	delete(p.hids, info.Hid)
	for _, it := range info.Hosts {
		delete(p.hosts, it)
	}
}

func (p *Connections) Bind(conn net.Conn, chs *Channels, host string, passive bool, fshaked func()) {
	rpc := NewChannels(chs.Sub("inf"))

	setcurr := func(info *ConnInfo, host string) {
		for i, it := range info.Hosts {
			if it == host {
				info.ConnHost = i
				return
			}
		}
		p.log.Msg("connected host should be listed")
		info.Hosts = append(info.Hosts, host)
		info.ConnHost = len(info.Hosts) - 1
	}

	add := func(node *NodeInfo) {
		p.lock.Lock()
		defer p.lock.Unlock()
		info := &ConnInfo{node, -1, conn, chs, nil, nil}
		setcurr(info, host)
		p.mconns[chs] = info
		p.hids[info.Hid] = info
		for _, it := range node.Hosts {
			if _, ok := p.hosts[it]; ok {
				p.log.Msg("error: reduplicated host")
			}
			p.hosts[it] = info
		}
		fshaked()
	}

	rpc.Sub(_ConnInfoHi).Receive(func(r io.Reader, n uint32) {
		info := &NodeInfo{}
		info.Load(r)
		add(info)
		data := p.info.Pack()
		rpc.Sub(_ConnInfoHiToo).Send(bytes.NewReader(data), uint32(len(data)))
	})

	rpc.Sub(_ConnInfoHiToo).Receive(func(r io.Reader, n uint32) {
		info := &NodeInfo{}
		info.Load(r)
		add(info)
	})

	if !passive {
		data := p.info.Pack()
		rpc.Sub(_ConnInfoHi).Send(bytes.NewReader(data), uint32(len(data)))
	}
}

func NewConnections(log *tools.Log, info * NodeInfo) *Connections {
	return &Connections{
		log,
		info,
		make(map[*Channels]*ConnInfo),
		make(map[*Channels]*ConnInfo),
		make(map[uint64]*ConnInfo),
		make(map[string]*ConnInfo),
		sync.Mutex{}}
}

type Connections struct {
	log *tools.Log
	info *NodeInfo
	mconns map[*Channels]*ConnInfo
	tconns map[*Channels]*ConnInfo
	hids map[uint64]*ConnInfo
	hosts map[string]*ConnInfo
	lock sync.Mutex
}

type ConnInfo struct {
	*NodeInfo
	ConnHost int
	MainConn net.Conn
	MainChs *Channels
	TransConn net.Conn
	TransChs *Channels
}

const (
	_ConnInfoHi = "hi"
	_ConnInfoHiToo = "hitoo"
)

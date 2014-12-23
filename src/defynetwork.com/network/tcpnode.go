package network

import (
	"net"
	"strings"
	"sync"
	"defynetwork.com/tools"
)

func (p *TcpNode) Serve() {
	p.svr.Serve(
		p.port,
		p.preread,
		func(conn net.Conn, chs *Channels) {
			host := strings.Split(conn.RemoteAddr().String(), ":")[0]
			p.lock.Lock()
			delete(p.connecting, host)
			p.connected[host] = true
			p.lock.Unlock()
			p.fconnected(conn, chs, true)
		},
		func(conn net.Conn, chs *Channels) {
			host := strings.Split(conn.RemoteAddr().String(), ":")[0]
			p.lock.Lock()
			delete(p.connected, host)
			p.lock.Unlock()
			p.fclosed(conn, chs)
		})
}

func (p *TcpNode) Conn(host string, port int) {
	for _, it := range p.addrs.MyAddrs(true) {
		if it == host && port == p.port {
			return
		}
	}

	p.lock.Lock()
	_, ok1 := p.connecting[host]
	_, ok2 := p.connected[host]
	if ok1 || ok2 {
		p.lock.Unlock()
		return
	}
	p.connecting[host] = true
	p.lock.Unlock()

	defer func() {
		err := recover()
		if err != nil {
			p.lock.Lock()
			delete(p.connecting, host)
			p.lock.Unlock()
			panic(err)
		}
	}()

	NewTcpClient(p.log).Start(
		host,
		port,
		p.preread,
		func(conn net.Conn, chs *Channels) {
			host := strings.Split(conn.RemoteAddr().String(), ":")[0]
			p.lock.Lock()
			delete(p.connecting, host)
			p.connected[host] = true
			p.lock.Unlock()
			p.fconnected(conn, chs, false)
		},
		func(conn net.Conn, chs *Channels) {
			host := strings.Split(conn.RemoteAddr().String(), ":")[0]
			p.lock.Lock()
			delete(p.connected, host)
			p.lock.Unlock()
			p.fclosed(conn, chs)
		})
}

func NewTcpNode(log *tools.Log, addrs *MyAddrs, port int, preread bool, fconnected ConnFuncEx, fclosed ConnFunc) *TcpNode {
	return &TcpNode{
		log,
		NewTcpSvr(log),
		addrs,
		port,
		preread,
		fconnected,
		fclosed,
		make(map[string]bool),
		make(map[string]bool),
		sync.Mutex{},
	}
}

type TcpNode struct {
	log *tools.Log
	svr *TcpSvr
	addrs *MyAddrs
	port int
	preread bool
	fconnected ConnFuncEx
	fclosed ConnFunc
	connected map[string]bool
	connecting map[string]bool
	lock sync.Mutex
}

type ConnFuncEx func(conn net.Conn, chs *Channels, passive bool)

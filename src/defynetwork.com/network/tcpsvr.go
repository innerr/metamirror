package network

import (
	"net"
	"strconv"
	"defynetwork.com/tools"
)

func (p *TcpSvr) Serve(port int, preread bool, fconnected, fclosed ConnFunc) {
	ln, err := net.Listen("tcp", ":" + strconv.Itoa(port))
	if err != nil {
		panic(err)
	}
	p.log.Msg("listen: " + strconv.Itoa(port))

	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		go func(conn net.Conn) {
			p.log.Msg("accept: " + conn.RemoteAddr().String())
			ch := NewTcpChannel(p.log, conn, preread)
			chs := NewChannels(ch)
			ch.OnClose(func(conn net.Conn) {
				p.log.Msg("close: " + conn.RemoteAddr().String())
				if fclosed != nil {
					fclosed(conn, chs)
				}
			})
			if fconnected != nil {
				fconnected(conn, chs)
			}
			ch.Start()
		}(conn)
	}
}

func NewTcpSvr(log *tools.Log) *TcpSvr {
	return &TcpSvr{log}
}

type TcpSvr struct {
	log *tools.Log
}

func (p *TcpClient) Start(host string, port int, preread bool, fconnected, fclosed ConnFunc) {
	conn, err := net.Dial("tcp", host + ":" + strconv.Itoa(port))
	if err != nil {
		panic(err)
	}
	p.log.Msg("conn: " + conn.RemoteAddr().String())
	ch := NewTcpChannel(p.log, conn, preread)
	chs := NewChannels(ch)
	if fconnected != nil {
		fconnected(conn, chs)
	}
	ch.OnClose(func(conn net.Conn) {
		p.log.Msg("deconn: " + conn.RemoteAddr().String())
		if fclosed != nil {
			fclosed(conn, chs)
		}
	})
	ch.Start()
}

func NewTcpClient(log *tools.Log) *TcpClient {
	return &TcpClient{log}
}

type TcpClient struct {
	log *tools.Log
}

type ConnFunc func(net.Conn, *Channels)

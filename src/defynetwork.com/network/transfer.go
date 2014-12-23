package network

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"net"
	"io"
	"strings"
	"time"
	"defynetwork.com/tools"
)

func (p *Transfer) TaskCount() int {
	return p.reqs.Count() + p.outs.Count()
}

func (p *Transfer) ReqCount() int {
	return p.reqs.Count()
}

func (p *Transfer) OutCount() int {
	return p.outs.Count()
}

func (p *Transfer) Write(hash []byte, hid uint64, fun func(r io.Reader, n uint32)) (done bool) {
	esha1 := hex.EncodeToString(hash)
	p.log.Debug("write begin: ", esha1)

	if bytes.Equal(hash, sha1.New().Sum(nil)) {
		fun(new(bytes.Buffer), 0)
		return true
	}

	if file, size := p.hashs.Get(hash); file != nil {
		defer file.Close()
		fun(file, size)
		p.log.Debug("from local: ", esha1)
		return true
	}
	p.reqs.Add(esha1, hash, fun)
	p.request(hash, esha1, hid)
	return false
}

func (p *Transfer) Conn(info *ConnInfo) {
	if info == nil {
		p.log.Debug("warn: null conn info")
		return
	}
	host := info.Hosts[info.ConnHost]
	go func() {
		defer tools.Catch(p.log.Detail, true)
		p.TcpNode.Conn(host, info.Tport)
	}()
}

func (p *Transfer) conn(hid uint64) {
	info := p.conns.ByHid(hid)
	if info == nil {
		p.log.Debug("can't find: ", fmt.Sprintf("%x", hid))
		return
	}
	p.Conn(info)
}

func (p *Transfer) LoopRetry(interval int) {
	for _ = range time.Tick(time.Duration(interval) * time.Millisecond) {
		for _, chs := range p.conns.TransConns() {
			p.refetch(chs)
		}
	}
}

func (p *Transfer) refetch(chs *Channels) {
	p.reqs.Walk(func(req *Request) {
		data := tools.Packd(req.Hash)
		chs.Sub(_TransferDataReq).Send(bytes.NewReader(data), uint32(len(data)))
		req.Rtime = time.Now().UnixNano()
	})
}

func (p Transfer) send(ch IChannel, data []byte) {
	go func() {
		ch.Send(bytes.NewReader(data), uint32(len(data)))
	}()
}

func (p *Transfer) request(hash []byte, esha1 string, hid uint64) {
	p.log.Debug("hash req: ", esha1)
	info := p.conns.ByHid(hid)
	if info != nil && info.TransChs != nil {
		p.send(info.TransChs.Sub(_TransferDataReq), tools.Packd(hash))
		p.reqs.Requesting(esha1)
	} else {
		p.conn(hid)
	}
}

func (p *Transfer) bind(conn net.Conn, rpc *Channels) {
	rpc.Sub(_TransferHashReq).Receive(func(r io.Reader, n uint32) {
		hash := tools.Loadd(r)
		esha1 := hex.EncodeToString(hash)
		p.log.Debug("hash req recv: ", esha1)
		if file, _ := p.hashs.Get(hash); file == nil {
			p.log.Debug("hash req recv: ", esha1, ", not found")
			return
		} else {
			file.Close()
		}
		p.log.Debug("hash resp: ", esha1)
		p.send(rpc.Sub(_TransferHashResp), tools.Packd(hash))
	})

	rpc.Sub(_TransferHashResp).Receive(func(r io.Reader, n uint32) {
		hash := tools.Loadd(r)
		esha1 := hex.EncodeToString(hash)
		found := p.reqs.Found(esha1, rpc)
		if found <= 0 {
			p.log.Debug("hash resp recv: ", esha1, ", not request")
			return
		}
		p.log.Debug("hash resp recv: ", esha1, " found:", found)
		p.log.Debug("data req: ", esha1)
		p.send(rpc.Sub(_TransferDataReq), tools.Packd(hash))
	})

	rpc.Sub(_TransferDataReq).Receive(func(r io.Reader, n uint32) {
		hash := tools.Loadd(r)
		esha1 := hex.EncodeToString(hash)
		file, size := p.hashs.Get(hash)
		if file == nil {
			p.log.Debug("data req recv: ", esha1, ", not found")
			return
		}
		defer file.Close()

		p.log.Debug("data req recv: ", esha1)
		p.log.Debug("data resp: ", esha1)

		go func() {
			defer func() {
				file.Close()
				p.outs.Done(hash)
			}()
			p.outs.Add(hash)
			resp := tools.Packd(hash)
			rpc.Sub(_TransferDataResp).Send(io.MultiReader(bytes.NewReader(resp), file), uint32(len(resp)) + uint32(size))
		}()
	})

	rpc.Sub(_TransferDataResp).Receive(func(r io.Reader, n uint32) {
		hash := tools.Loadd(r)
		esha1 := hex.EncodeToString(hash)
		size := n - uint32(len(tools.Packd(hash)))
		total, done := p.reqs.Done(esha1, r, size)
		if total == 0 {
			tools.Eat(p.log, r, size)
		}
		p.log.Debug("data resp recv: ", esha1, " done: ", done, "/", total)
	})
}

func NewTransfer(log *tools.Log, addrs *MyAddrs, port int, hashs IHashLib, conns *Connections, timeout int) *Transfer {
	p := &Transfer{log, nil, hashs, conns, NewRequests(log.Mod("transreqs"), timeout), NewTransOuts(log.Mod("transouts"))}
	p.TcpNode = NewTcpNode(
		log,
		addrs,
		port,
		false,
		func(conn net.Conn, chs *Channels, passive bool) {
			host := strings.Split(conn.RemoteAddr().String(), ":")[0]
			conns.AddTransConn(conn, chs, host)
			p.bind(conn, chs)
			p.refetch(chs)
		},
		func(conn net.Conn, chs *Channels) {
			p.conns.DelTransConn(chs)
			chs.CleanAll()
		})
	return p
}

type Transfer struct {
	log *tools.Log
	*TcpNode
	hashs IHashLib
	conns *Connections
	reqs *Requests
	outs *TransOuts
}

const (
	_TransferHashReq = "hashreq"
	_TransferHashResp = "hashresp"
	_TransferDataReq = "datareq"
	_TransferDataResp = "dataresp"
)

type IHashLib interface {
	Get([]byte) (io.ReadCloser, uint32)
}

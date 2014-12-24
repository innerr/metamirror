package structs

import (
	"bytes"
	"io"
	"sync"
	"defynetwork.com/tools"
)

func (p *Session) Sync() {
	if p.core.Flags.In {
		p.cliSendSyncReq()
	} else {
		p.cliSendRoReq()
	}
}

func (p *Session) cliSendSyncReq() {
	p.log.Debug("sync")
	buf := new(bytes.Buffer)
	p.core.Clocks().Dump(buf)
	tools.Dumpb(buf, p.core.Flags.Integral)
	data := buf.Bytes()
	p.rpc.Func(_SyncReq).Send(bytes.NewReader(data), uint32(len(data)))
}

func (p *Session) svrRecvSyncReq(r io.Reader, n uint32) {
	if !p.core.Flags.Out || p.blocked {
		p.svrSendRoResp()
		return
	}

	c := Clocks{}
	c.Load(r)
	integral := tools.Loadb(r)
	var delta Delta
	var clocks Clocks
	if integral {
		p.log.Debug("syncreq recv: integral")
		delta, clocks = p.core.Pack()
	} else {
		p.log.Debug("syncreq recv")
		delta, clocks = p.core.Delta(c)
	}
	p.svrSendSyncResp(delta, clocks)
}

func (p *Session) svrSendSyncResp(delta Delta, clocks Clocks) {
	buf := new(bytes.Buffer)
	clocks.Dump(buf)
	delta.Dump(buf)
	data := buf.Bytes()
	p.log.Debug("syncresp")
	p.rpc.Func(_SyncResp).Send(bytes.NewReader(data), uint32(len(data)))
	p.synced = true
}

func (p *Session) cliRecvSyncResp(r io.Reader, n uint32) {
	p.log.Debug("syncresp recv")
	c := Clocks{}
	c.Load(r)
	d := Delta{}
	d.Load(r)
	p.core.Merge(d)
	delta, _ := p.core.Delta(c)
	p.cliSendComplete(delta)
}

func (p *Session) cliSendComplete(delta Delta) {
	if !p.core.Flags.Out || p.blocked {
		return
	}
	buf := new(bytes.Buffer)
	delta.Dump(buf)
	data := buf.Bytes()
	if len(delta) != 0 {
		p.log.Debug("complete")
		p.rpc.Func(_SyncComplete).Send(bytes.NewReader(data), uint32(len(data)))
	} else {
		p.log.Debug("complete: skip")
	}
	p.synced = true
}

func (p *Session) svrRecvComplete(r io.Reader, n uint32) {
	p.log.Debug("complete recv")
	d := Delta{}
	d.Load(r)
	p.core.Merge(d)
}

func (p *Session) cliSendRoReq() {
	p.log.Debug("roreq")
	p.rpc.Func(_SyncRoReq).Send(bytes.NewReader([]byte{}), 0)
}

func (p *Session) svrRecvRoReq(r io.Reader, n uint32) {
	p.log.Debug("roreq recv")
	p.blocked = true
	p.svrSendRoResp()
}

func (p *Session) svrSendRoResp() {
	buf := new(bytes.Buffer)
	p.core.Clocks().Dump(buf)
	data := buf.Bytes()
	p.log.Debug("roresp")
	p.rpc.Func(_SyncRoResp).Send(bytes.NewReader(data), uint32(len(data)))
}

func (p *Session) cliRecvRoResp(r io.Reader, n uint32) {
	p.log.Debug("roresp recv")
	c := Clocks{}
	c.Load(r)
	delta, _ := p.core.Delta(c)
	p.cliSendComplete(delta)
}

func (p *Session) SendDelta(hid uint64, delta Delta) {
	if !p.core.Flags.Broadcast || p.blocked {
		return
	}

	p.log.Debug("send delta")
	buf := new(bytes.Buffer)
	tools.Dump(buf, hid)
	delta.Dump(buf)
	data := buf.Bytes()
	p.rpc.Func(_SyncDelta).Send(bytes.NewReader(data), uint32(len(data)))
}

func (p *Session) recvDelta(r io.Reader, n uint32) {
	p.log.Debug("recv delta")

	hid := tools.Loadu64(r)
	d := Delta{}
	d.Load(r)

	if p.frecv != nil {
		p.frecv(p.ch, hid, d)
	}

	if !p.core.Flags.ManualMerge {
		p.core.Merge(d)
		return
	}

	p.lock.Lock()
	defer p.lock.Unlock()
	p.received = append(p.received, d)
}

func (p *Session) ManualMerge() {
	p.lock.Lock()
	defer p.lock.Unlock()
	for _, v := range p.received {
		p.core.Merge(v)
	}
	p.received = nil
}

func (p *Session) Close() {
	p.ManualMerge()
	p.rpc.Close()
}

func NewSession(core *Core, ch IChannel, frecv func(IChannel, uint64, Delta), log *tools.Log) *Session {
	p := &Session{
		core: core,
		ch: ch,
		rpc: NewRpc(ch),
		frecv: frecv,
		log: log,
	}

	p.rpc.Func(_SyncDelta).Receive(p.recvDelta)

	p.rpc.Func(_SyncReq).Receive(p.svrRecvSyncReq)
	p.rpc.Func(_SyncResp).Receive(p.cliRecvSyncResp)
	p.rpc.Func(_SyncComplete).Receive(p.svrRecvComplete)

	p.rpc.Func(_SyncRoReq).Receive(p.svrRecvRoReq)
	p.rpc.Func(_SyncRoResp).Receive(p.cliRecvRoResp)

	return p
}

type Session struct {
	core *Core
	ch IChannel
	rpc *Rpc
	synced bool
	blocked bool
	received []Delta
	frecv func(IChannel, uint64, Delta)
	lock sync.Mutex
	log *tools.Log
}

const (
	_SyncDelta = "delta"
	_SyncReq = "syncreq"
	_SyncResp = "syncresp"
	_SyncComplete = "complete"
	_SyncRoReq = "roreq"
	_SyncRoResp = "roresp"
)

package structs

import (
	"bytes"
	"io"
	"sync"
	"defynetwork.com/tools"
)

func (p *Session) MergeReceived() {
	p.lock.Lock()
	defer p.lock.Unlock()
	for _, v := range p.received {
		p.core.Merge(v)
	}
	p.received = nil
}

func (p *Session) Received(fun func()) {
	p.lock.Lock()
	defer p.lock.Unlock()
	if p.frecv != nil {
		panic("double assigned")
	}
	p.frecv = fun
}

func (p *Session) syncreq(rpc *Channels, clocks Clocks) {
	p.log.Debug("sync")
	buf := new(bytes.Buffer)
	clocks.Dump(buf)
	tools.Dumpb(buf, p.core.Flags.Integral)
	data := buf.Bytes()
	rpc.Sub(_SyncReq).Send(bytes.NewReader(data), uint32(len(data)))
}

func (p *Session) syncresp(rpc *Channels, clocks Clocks, delta Delta) {
	p.log.Debug("syncresp")
	buf := new(bytes.Buffer)
	clocks.Dump(buf)
	delta.Dump(buf)
	data := buf.Bytes()
	rpc.Sub(_SyncResp).Send(bytes.NewReader(data), uint32(len(data)))
	p.synced = true
}

func (p *Session) complete(rpc *Channels, delta Delta) {
	if !p.core.Flags.Out || p.blocked != false {
		return
	}
	buf := new(bytes.Buffer)
	delta.Dump(buf)
	data := buf.Bytes()
	if len(delta) != 0 {
		p.log.Debug("fsync.complete")
		rpc.Sub(_SyncComplete).Send(bytes.NewReader(data), uint32(len(data)))
	} else {
		p.log.Debug("fsync.complete: skip")
	}
	p.synced = true
}

rpc.Sub(_SyncDelta).Receive(func(r io.Reader, n uint32) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.log.Debug("recv delta")
	hid := tools.Loadu64(r)
	d := Delta{}
	d.Load(r)
	if p.core.Flags.ManualMerge {
		p.received = append(p.received, d)
		if p.frecv != nil {
			p.frecv()
		}
	} else {
		p.core.Merge(d)
	}
	if !p.core.Flags.Broadcast {
		return
	}
	if len(p.conns) <= 1 {
		return
	}
	p.log.Debug("resend to: ", len(p.conns), "-1")
	for _, v := range p.conns {
		if v != rpc && p.synced != false {
			p.send(rpc, hid, d)
		}
	}
})

rpc.Sub(_SyncReq).Receive(func(r io.Reader, n uint32) {
	p.lock.Lock()
	defer p.lock.Unlock()
	c := Clocks{}
	c.Load(r)
	integral := tools.Loadb(r)
	var delta Delta
	if integral {
		p.log.Debug("fsync.syncreq recv: integral")
		delta = p.core.Pack()
	} else {
		p.log.Debug("fsync.syncreq recv")
		delta = p.core.Delta(c)
	}
	syncresp(rpc, p.core.Clocks(), delta)
})

rpc.Sub(_SyncResp).Receive(func(r io.Reader, n uint32) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.log.Debug("fsync.syncresp recv")
	c := Clocks{}
	c.Load(r)
	d := Delta{}
	d.Load(r)
	p.core.Merge(d)
	delta := p.core.Delta(c)
	complete(rpc, delta)
})

rpc.Sub(_SyncComplete).Receive(func(r io.Reader, n uint32) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.log.Debug("fsync.complete recv")
	d := Delta{}
	d.Load(r)
	p.core.Merge(d)
})

rpc.Sub(_SyncRoReq).Receive(func(r io.Reader, n uint32) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.log.Debug("fsync.roreq recv")
	buf := new(bytes.Buffer)
	p.core.Clocks().Dump(buf)
	data := buf.Bytes()
	p.blocked = true
	p.log.Debug("fsync.roresp")
	rpc.Sub(_SyncRoResp).Send(bytes.NewReader(data), uint32(len(data)))
})

rpc.Sub(_SyncRoResp).Receive(func(r io.Reader, n uint32) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.log.Debug("fsync.roresp recv")
	c := Clocks{}
	c.Load(r)
	delta := p.core.Delta(c)
	p.log.Debug("fsync.complete")
	complete(rpc, delta)
	p.synced = true
})

p.conns[ch] = rpc
if !passive {
	if p.core.Flags.In {
		syncreq(rpc, p.core.Clocks())
	} else {
		rpc.Sub(_SyncRoReq).Send(bytes.NewReader([]byte{}), uint32(0))
	}
}

func NewSession(log *tools.Log, core *Core, ch IChannal) *Session {
	p := &Session{
		log: log,
		core: core,
		ch: ch,
		rpc: NewChannels(ch),
	}
	return p
}

type Session struct {
	log *tools.Log
	core *Core
	ch IChannal
	rpc *Channals
	synced bool
	blocked bool
	received []Delta
	frecv func()
	lock sync.Mutex
}

const (
	_SyncDelta = "delta"
	_SyncReq = "syncreq"
	_SyncResp = "syncresp"
	_SyncComplete = "complete"
	_SyncRoReq = "roreq"
	_SyncRoResp = "roresp"
)

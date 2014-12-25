package structs

import (
	"bytes"
	"io"
	"sync"
	"defynetwork.com/tools"
)

func (p *Session) Sync() {
	p.cliSendSyncReq()
}

func (p *Session) cliSendSyncReq() {
	buf := new(bytes.Buffer)
	p.core.Clocks().Dump(buf)
	tools.Dumpb(buf, p.core.Flags.In)
	tools.Dumpb(buf, p.core.Flags.Integral)
	data := buf.Bytes()
	p.log.Debug("req")
	p.rpc.Func(_SyncRequest).Send(bytes.NewReader(data), uint32(len(data)))
}

func (p *Session) svrRecvSyncReq(r io.Reader, n uint32) {
	clocks := Clocks{}
	clocks.Load(r)
	p.blocked = !tools.Loadb(r)
	integral := tools.Loadb(r)
	if integral {
		p.log.Debug("req recv: integral")
	} else {
		p.log.Debug("req recv")
	}
	p.svrSendSyncResp(integral, clocks)
}

func (p *Session) svrSendSyncResp(integral bool, remote Clocks) {
	buf := new(bytes.Buffer)
	tools.Dumpb(buf, p.core.Flags.In)

	if p.core.Flags.Out && !p.blocked {
		var delta Delta
		var clocks Clocks
		if integral {
			delta, clocks = p.core.Pack()
		} else {
			delta, clocks = p.core.Delta(remote)
		}
		clocks.Dump(buf)
		delta.Dump(buf)
		p.log.Debug("resp")
	} else {
		p.log.Debug("resp: !out|!cli.in")
	}

	data := buf.Bytes()
	p.rpc.Func(_SyncResponse).Send(bytes.NewReader(data), uint32(len(data)))
	p.setSended()
}

func (p *Session) cliRecvSyncResp(r io.Reader, n uint32) {
	p.log.Debug("resp recv")
	p.blocked = !tools.Loadb(r)

	if n > tools.DumpbSize {
		clocks := Clocks{}
		clocks.Load(r)
		delta := Delta{}
		delta.Load(r)
		p.core.Merge(delta)
		p.cliSendComplete(clocks)
	} else {
		p.cliSendComplete(nil)
	}
}

func (p *Session) cliSendComplete(remote Clocks) {
	delta, _ := p.core.Delta(remote)

	if !p.core.Flags.Out || p.blocked {
		p.log.Debug("complete: !out|!cli.in")
		p.rpc.Func(_SyncComplete).Send(bytes.NewReader([]byte{}), 0)
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
		p.rpc.Func(_SyncComplete).Send(bytes.NewReader([]byte{}), 0)
	}
	p.setSended()
}

func (p *Session) svrRecvComplete(r io.Reader, n uint32) {
	p.log.Debug("complete recv")
	if n != 0 {
		d := Delta{}
		d.Load(r)
		p.core.Merge(d)
	}
	p.setSynced()
	p.svrSendConfirm()
}

func (p *Session) svrSendConfirm() {
	p.log.Debug("confirm")
	p.rpc.Func(_SyncConfirm).Send(bytes.NewReader([]byte{}), 0)
}

func (p *Session) cliRecvConfirm(r io.Reader, n uint32) {
	p.setSynced()
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

func (p *Session) setSended() {
	p.log.Debug("sended")
	p.sended = true
	if !p.core.Flags.In {
		p.setSynced()
	}
}

func (p *Session) setSynced() {
	p.log.Debug("synced")
	p.synced = true
	if p.fsynced != nil {
		p.fsynced()
	}
}

func NewSession(core *Core, ch IChannel, frecv RecvDeltaFunc, fsynced FullSyncedFunc, log *tools.Log) *Session {
	p := &Session{
		core: core,
		ch: ch,
		rpc: NewRpc(ch),
		frecv: frecv,
		fsynced: fsynced,
		log: log,
	}

	p.rpc.Func(_SyncDelta).Receive(p.recvDelta)

	p.rpc.Func(_SyncRequest).Receive(p.svrRecvSyncReq)
	p.rpc.Func(_SyncResponse).Receive(p.cliRecvSyncResp)
	p.rpc.Func(_SyncComplete).Receive(p.svrRecvComplete)
	p.rpc.Func(_SyncConfirm).Receive(p.cliRecvConfirm)

	return p
}

type Session struct {
	core *Core
	ch IChannel
	rpc *Rpc
	received []Delta
	frecv RecvDeltaFunc
	fsynced FullSyncedFunc
	sended bool
	synced bool
	blocked bool
	lock sync.Mutex
	log *tools.Log
}

const (
	_SyncDelta = "delta"
	_SyncRequest = "request"
	_SyncResponse = "response"
	_SyncComplete = "complete"
	_SyncConfirm = "confirm"
)

type FullSyncedFunc func()
type RecvDeltaFunc func(IChannel, uint64, Delta)

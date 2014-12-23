package structs

import (
	"bytes"
	"io"
	"sync"
	"defynetwork.com/tools"
)

// TODO: remove closure functions

func (p *Domain) MergeReceived() {
	p.lock.Lock()
	defer p.lock.Unlock()
	for _, v := range p.received {
		p.merge(v)
	}
	p.received = nil
}

func (p *Domain) Received(fun func()) {
	p.lock.Lock()
	defer p.lock.Unlock()
	if p.frecv != nil {
		panic("double assigned")
	}
	p.frecv = fun
}

func (p *Domain) Clocks() Clocks {
	p.lock.Lock()
	defer p.lock.Unlock()
	return p.box.Max()
}

func (p *Domain) Delta(clocks Clocks) Delta {
	p.lock.Lock()
	defer p.lock.Unlock()
	return p.box.Delta(clocks)
}

func (p *Domain) Datas() *Compose {
	return p.data
}

func (p *Domain) Merge(delta Delta) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.merge(delta)
}

func (p *Domain) Bind(ch IChannel, passive bool) bool {
	p.lock.Lock()
	defer p.lock.Unlock()

	if _, ok := p.conns[ch]; ok {
		return false
	}
	p.log.Debug("bind ", passive)
	rpc := NewChannels(ch)

	syncreq := func(rpc *Channels, clocks Clocks) {
		p.log.Debug("fsync.sync")
		buf := new(bytes.Buffer)
		clocks.Dump(buf)
		tools.Dumpb(buf, p.flags.Integral)
		data := buf.Bytes()
		rpc.Sub(_SyncReq).Send(bytes.NewReader(data), uint32(len(data)))
	}

	syncresp := func(rpc *Channels, clocks Clocks, delta Delta) {
		p.log.Debug("fsync.syncresp")
		buf := new(bytes.Buffer)
		clocks.Dump(buf)
		delta.Dump(buf)
		data := buf.Bytes()
		rpc.Sub(_SyncResp).Send(bytes.NewReader(data), uint32(len(data)))
		p.synced[rpc] = true
	}

	complete := func(rpc *Channels, delta Delta) {
		if !p.flags.Out || p.blocked[rpc] != false {
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
		p.synced[rpc] = true
	}

	rpc.Sub(_SyncDelta).Receive(func(r io.Reader, n uint32) {
		p.lock.Lock()
		defer p.lock.Unlock()
		p.log.Debug("recv delta")
		hid := tools.Loadu64(r)
		d := Delta{}
		d.Load(r)
		if p.flags.ManualMerge {
			p.received = append(p.received, d)
			if p.frecv != nil {
				p.frecv()
			}
		} else {
			p.merge(d)
		}
		if !p.flags.Broadcast {
			return
		}
		if len(p.conns) <= 1 {
			return
		}
		p.log.Debug("resend to: ", len(p.conns), "-1")
		for _, v := range p.conns {
			if v != rpc && p.synced[v] != false {
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
			clocks := p.box.Max()
			delta = Delta{p.data.Pack(clocks)}
		} else {
			p.log.Debug("fsync.syncreq recv")
			delta = p.box.Delta(c)
		}
		syncresp(rpc, p.box.Max(), delta)
	})

	rpc.Sub(_SyncResp).Receive(func(r io.Reader, n uint32) {
		p.lock.Lock()
		defer p.lock.Unlock()
		p.log.Debug("fsync.syncresp recv")
		c := Clocks{}
		c.Load(r)
		d := Delta{}
		d.Load(r)
		p.merge(d)
		delta := p.box.Delta(c)
		complete(rpc, delta)
	})

	rpc.Sub(_SyncComplete).Receive(func(r io.Reader, n uint32) {
		p.lock.Lock()
		defer p.lock.Unlock()
		p.log.Debug("fsync.complete recv")
		d := Delta{}
		d.Load(r)
		for _, blob := range d {
			p.persist.Dump(blob)
		}
		p.merge(d)
	})

	rpc.Sub(_SyncRoReq).Receive(func(r io.Reader, n uint32) {
		p.lock.Lock()
		defer p.lock.Unlock()
		p.log.Debug("fsync.roreq recv")
		buf := new(bytes.Buffer)
		p.box.Max().Dump(buf)
		data := buf.Bytes()
		p.blocked[rpc] = true
		p.log.Debug("fsync.roresp")
		rpc.Sub(_SyncRoResp).Send(bytes.NewReader(data), uint32(len(data)))
	})

	rpc.Sub(_SyncRoResp).Receive(func(r io.Reader, n uint32) {
		p.lock.Lock()
		defer p.lock.Unlock()
		p.log.Debug("fsync.roresp recv")
		c := Clocks{}
		c.Load(r)
		delta := p.box.Delta(c)
		p.log.Debug("fsync.complete")
		complete(rpc, delta)
		p.synced[rpc] = true
	})

	p.conns[ch] = rpc
	if !passive {
		if p.flags.In {
			syncreq(rpc, p.box.Max())
		} else {
			rpc.Sub(_SyncRoReq).Send(bytes.NewReader([]byte{}), uint32(0))
		}
	}
	return true
}

func (p *Domain) Unbind(ch IChannel) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.log.Debug("unbind")
	delete(p.conns, ch)
	ch.Receive(nil)
}

func (p *Domain) Commit(hid uint64) {
	p.lock.Lock()
	defer p.lock.Unlock()

	clocks := p.box.Edit(hid)
	blob := p.data.Commit(clocks)
	if blob.IsNil() {
		return
	}

	p.log.Debug("commit to: ", len(p.conns))
	p.box.Add(blob, hid)
	p.persist.Dump(blob)
	d := Delta{blob}

	for _, v := range p.conns {
		p.send(v, hid, d)
	}
}

func (p *Domain) merge(delta Delta) {
	p.log.Debug("merge")
	for _, blob := range delta {
		p.persist.Dump(blob)
	}
	p.box.Merge(delta)
	for _, blob := range delta {
		p.data.Merge(blob)
	}
}

func (p *Domain) send(rpc *Channels, hid uint64, delta Delta) {
	if !p.flags.Out || p.blocked[rpc] != false {
		return
	}
	p.log.Debug("send delta")
	buf := new(bytes.Buffer)
	tools.Dump(buf, hid)
	delta.Dump(buf)
	data := buf.Bytes()
	rpc.Sub(_SyncDelta).Send(bytes.NewReader(data), uint32(len(data)))
}

func NewDomain(log *tools.Log, data IData, persist IPersist, flags *DomainFlags) *Domain {
	if flags == nil {
		flags = NewDomainFlags()
	}
	p := &Domain{
		log,
		NewCompose(data),
		NewBox(),
		persist,
		flags,
		make(map[IChannel]*Channels),
		make(map[*Channels]bool),
		make(map[*Channels]bool),
		nil,
		nil,
		sync.Mutex{},
	}
	p.persist.Load(func(blob Blob) {
		p.box.Add(blob, blob.Vcs.Max())
		p.data.Merge(blob)
	})
	return p
}

type Domain struct {
	log *tools.Log
	data *Compose
	box *Box
	persist IPersist
	flags *CoreFlags
	conns map[IChannel]*Channels
	synced map[*Channels]bool
	blocked map[*Channels]bool
	received []Delta
	frecv func()
	lock sync.Mutex
}

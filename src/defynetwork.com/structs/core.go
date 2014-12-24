package structs

import (
	"sync"
	"defynetwork.com/tools"
)

func (p *Core) Clocks() Clocks {
	p.lock.Lock()
	defer p.lock.Unlock()
	return p.box.Max()
}

func (p *Core) Pack() (Delta, Clocks) {
	p.lock.Lock()
	defer p.lock.Unlock()
	clocks := p.box.Max()
	blob := Blob{nil, nil}
	data := p.data.Pack()
	if len(data) != 0 {
		blob = Blob{clocks, data}
	}
	return Delta{blob}, clocks
}

func (p *Core) Delta(clocks Clocks) (Delta, Clocks) {
	p.lock.Lock()
	defer p.lock.Unlock()
	return p.box.Delta(clocks), p.box.Max()
}

func (p *Core) Merge(delta Delta) {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.log.Debug("merge")

	for _, blob := range delta {
		p.persist.Dump(blob)
	}
	p.box.Merge(delta)
	for _, blob := range delta {
		if len(blob.Data) != 0 {
			p.data.Merge(blob.Data)
		}
	}
}

func (p *Core) Commit(hid uint64) Delta {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.log.Debug("committing")

	clocks := p.box.Edit(hid)
	data := p.data.Commit()
	if len(data) == 0 {
		return nil
	}
	blob := Blob{clocks, data}

	p.persist.Dump(blob)
	p.box.Add(blob, hid)

	p.log.Debug("commited")
	return Delta{blob}
}

func (p *Core) load(blob Blob) {
	p.box.Add(blob, blob.Vcs.Max())
	if len(blob.Data) != 0 {
		p.data.Merge(blob.Data)
	}
}

func NewCore(data IData, persist IPersist, flags *CoreFlags, log *tools.Log) *Core {
	if flags == nil {
		flags = NewCoreFlags()
	}
	p := &Core{
		data: data,
		box: NewBox(),
		persist: persist,
		Flags: flags,
		log: log,
	}
	p.persist.Load(p.load)
	return p
}

type Core struct {
	data IData
	box *Box
	persist IPersist
	Flags *CoreFlags
	lock sync.Mutex
	log *tools.Log
}

func NewCoreFlags() *CoreFlags {
	return &CoreFlags{true, true, true, false, false}
}

type CoreFlags struct {
	In bool
	Out bool
	Broadcast bool
	Integral bool
	ManualMerge bool
}

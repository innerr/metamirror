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
	return Delta{p.data.Pack(clocks)}, clocks
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
		p.data.Merge(blob)
	}
}

func (p *Core) Commit(hid uint64) Delta {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.log.Debug("committing")

	clocks := p.box.Edit(hid)
	blob := p.data.Commit(clocks)
	if blob.IsNil() {
		return nil
	}

	p.persist.Dump(blob)
	p.box.Add(blob, hid)

	p.log.Debug("commited")
	return Delta{blob}
}

func NewCore(log *tools.Log, data IData, persist IPersist, flags *CoreFlags) *Core {
	if flags == nil {
		flags = NewCoreFlags()
	}
	p := &Core{
		log: log,
		data: NewCompose(data),
		box: NewBox(),
		persist: persist,
		Flags: flags,
	}
	p.persist.Load(func(blob Blob) {
		p.box.Add(blob, blob.Vcs.Max())
		p.data.Merge(blob)
	})
	return p
}

type Core struct {
	log *tools.Log
	data *Compose
	box *Box
	persist IPersist
	Flags *CoreFlags
	lock sync.Mutex
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

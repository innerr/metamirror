package structs

import (
	"bytes"
	"io"
	"sync"
	"defynetwork.com/tools"
)

func (p *Domains) Bind(ch IChannel, passive bool) {
	p.lock.Lock()
	defer p.lock.Unlock()

	rpc := NewChannels(ch)
	p.conns[ch] = rpc

	grant := func(names []string, passive bool) []string {
		granted := []string{}
		for _, name := range names {
			domain, ok := p.regs[name]
			if !ok && p.newdm != nil {
				domain = p.newdm(name)
				p.Reg(name, domain)
			}
			if domain != nil {
				binded := domain.Bind(rpc.Sub(name), passive)
				if !binded {
					continue
				}
				granted = append(granted, name)
			}
		}
		p.log.Debug("granted", granted)
		return granted
	}

	ungrant := func(names []string) {
		p.lock.Lock()
		defer p.lock.Unlock()
		for _, name := range names {
			domain, ok := p.regs[name]
			if !ok {
				continue
			}
			domain.Unbind(rpc.Sub(name))
		}
	}

	rpc.Sub(_DomainsReq).Receive(func(r io.Reader, n uint32) {
		p.log.Debug("req recv")
		p.lock.Lock()
		defer p.lock.Unlock()

		names := tools.Unpackss(r)
		p.remotes[ch] = names
		granted := grant(names, true)
		d1 := tools.Packss(granted)

		regs := []string{}
		for name, _ := range p.regs {
			regs = append(regs, name)
		}
		d2 := tools.Packss(regs)

		p.log.Debug("grant send")
		rpc.Sub(_DomainsGrant).Send(io.MultiReader(bytes.NewReader(d1), bytes.NewReader(d2)), uint32(len(d1) + len(d2)))
	})

	rpc.Sub(_DomainsGrant).Receive(func(r io.Reader, n uint32) {
		p.log.Debug("grant recv")
		p.lock.Lock()
		defer p.lock.Unlock()
		granted := tools.Unpackss(r)
		remotes := tools.Unpackss(r)
		p.remotes[ch] = remotes
		grant(granted, false)
	})

	rpc.Sub(_DomainsDel).Receive(func(r io.Reader, n uint32) {
		p.log.Debug("del recv")
		names := tools.Unpackss(r)
		ungrant(names)
	})

	if !passive {
		p.sync(rpc)
	}
	p.log.Msg("bind")
}

func (p *Domains) Unbind(ch IChannel) {
	p.lock.Lock()
	defer p.lock.Unlock()

	rpc, ok := p.conns[ch]
	if !ok {
		return
	}
	p.log.Msg("unbind")
	for name, domain := range p.regs {
		domain.Unbind(rpc.Sub(name))
	}
	delete(p.conns, ch)
}

func (p *Domains) Unreg(name string) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.log.Msg("unreg: " + name)
	domain, ok := p.regs[name]
	if ok {
		delete(p.regs, name)
		for _, rpc := range p.conns {
			domain.Unbind(rpc.Sub(name))
		}
	}
	p.unregs[name] = true
	p.modified = true
}

func (p *Domains) Reg(name string, domain *Domain) {
	if domain == nil {
		return
	}
	if _, ok := p.regs[name]; ok {
		panic("exists")
	}
	p.lock.Lock()
	defer p.lock.Unlock()
	p.log.Msg("reg: " + name)
	p.regs[name] = domain
	delete(p.unregs, name)
	p.modified = true
}

func (p *Domains) Sync() {
	p.lock.Lock()
	defer p.lock.Unlock()
	if p.modified == false {
		return
	}
	p.log.Msg("sync")
	for _, rpc := range p.conns {
		p.sync(rpc)
	}
	p.modified = false
}

func (p *Domains) sync(rpc *Channels) {
	regs := []string{}
	for name, _ := range p.regs {
		regs = append(regs, name)
	}
	rd := tools.Packss(regs)
	unregs := []string{}
	for name, _ := range p.unregs {
		unregs = append(unregs, name)
	}
	ud := tools.Packss(unregs)
	p.log.Debug("req/del send")
	rpc.Sub(_DomainsReq).Send(bytes.NewReader(rd), uint32(len(rd)))
	rpc.Sub(_DomainsDel).Send(bytes.NewReader(ud), uint32(len(ud)))
}

func (p *Domains) MergeReceived() {
	p.lock.Lock()
	defer p.lock.Unlock()
	for _, domain := range p.regs {
		domain.MergeReceived()
	}
}

func (p *Domains) Commit(hid uint64) {
	p.lock.Lock()
	defer p.lock.Unlock()
	for _, domain := range p.regs {
		domain.Commit(hid)
	}
}

func (p *Domains) Get(name string) *Domain {
	p.lock.Lock()
	defer p.lock.Unlock()
	return p.regs[name]
}

func (p *Domains) RemoteDomains(ch IChannel) []string {
	p.lock.Lock()
	defer p.lock.Unlock()
	return p.remotes[ch]
}

func (p *Domains) Names() []string {
	p.lock.Lock()
	defer p.lock.Unlock()
	names := []string{}
	for k, _ := range p.regs {
		names = append(names, k)
	}
	return names
}

func NewDomains(log *tools.Log, newdm NewDmFunc) *Domains {
	return &Domains{
		log,
		make(map[string]*Domain),
		make(map[string]bool),
		make(map[IChannel]*Channels),
		make(map[IChannel][]string),
		newdm,
		false,
		sync.Mutex{}}
}

type Domains struct {
	log *tools.Log
	regs map[string]*Domain
	unregs map[string]bool
	conns map[IChannel]*Channels
	remotes map[IChannel][]string
	newdm NewDmFunc
	modified bool
	lock sync.Mutex
}

const (
	_DomainsReq = "request"
	_DomainsGrant = "grant"
	_DomainsDel = "remove"
)

type NewDmFunc func(domain string)*Domain

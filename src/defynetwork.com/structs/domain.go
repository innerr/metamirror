package structs

import (
	"sync"
	"defynetwork.com/tools"
)

func (p *Domain) Bind(ch IChannel, passive bool) bool {
	creater := func() *Session {
		return NewSession(p.core, ch, p.recvDelta, nil, p.log.Mod("session"))
	}
	session := p.conns.Add(ch, creater)
	if session == nil {
		return false
	}

	p.log.Debug("bind", ch, passive)
	if !passive {
		session.Sync()
	}
	return true
}

func (p *Domain) Unbind(ch IChannel) {
	p.log.Debug("unbind", ch)
	session := p.conns.Remove(ch)
	session.Close()
}

func (p *Domain) Commit(hid uint64) {
	delta := p.core.Commit(hid)
	p.sendDelta(nil, hid, delta)
}

func (p *Domain) ManualMerge() {
	for _, session := range p.conns.Data() {
		session.ManualMerge()
	}
}

func (p *Domain) recvDelta(ch IChannel, hid uint64, delta Delta) {
	p.sendDelta(ch, hid, delta)
}

func (p *Domain) sendDelta(from IChannel, hid uint64, delta Delta) {
	for ch, session := range p.conns.Data() {
		if ch == from {
			continue
		}
		session.SendDelta(hid, delta)
	}
}

func NewDomain(data IData, persist IPersist, flags *CoreFlags, log *tools.Log) *Domain {
	p := &Domain{
		core: NewCore(data, persist, flags, log.Mod("core")),
		conns: NewDomainConns(),
		log: log,
	}
	return p
}

type Domain struct {
	core *Core
	conns *DomainConns
	log *tools.Log
}

func (p *DomainConns) Data() map[IChannel]*Session {
	p.lock.Lock()
	defer p.lock.Unlock()
	conns := make(map[IChannel]*Session)
	for k, v := range p.data {
		conns[k] = v
	}
	return conns
}

func (p *DomainConns) Add(ch IChannel, creater func() *Session) *Session {
	p.lock.Lock()
	defer p.lock.Unlock()
	session := p.data[ch]
	if session != nil {
		return nil
	}
	p.data[ch] = creater()
	return session
}

func (p *DomainConns) Remove(ch IChannel) *Session {
	p.lock.Lock()
	defer p.lock.Unlock()
	session := p.data[ch]
	if session == nil {
		return nil
	}
	delete(p.data, ch)
	return session
}

func NewDomainConns() *DomainConns {
	return &DomainConns{data: make(map[IChannel]*Session)}
}

type DomainConns struct {
	data map[IChannel]*Session
	lock sync.Mutex
}

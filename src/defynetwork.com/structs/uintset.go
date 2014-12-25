package structs

import (
	"bytes"
	"defynetwork.com/tools"
)

func (p *UintSet) Set(val uint32) *UintSet {
	_, ok := p.vals[val]
	if ok {
		return p
	}
	p.vals[val] = true
	p.adds = append(p.adds, val)
	return p
}

func (p *UintSet) Del(val uint32) {
	_, ok := p.vals[val]
	if !ok {
		return
	}
	delete(p.vals, val)
	p.dels = append(p.dels, val)
}

func (p *UintSet) Has(val uint32) bool {
	_, ok := p.vals[val]
	return ok
}

func (p *UintSet) Equal(x *UintSet) bool {
	if len(x.vals) != len(p.vals) {
		return false
	}
	for it, _ := range x.vals {
		if _, ok := p.vals[it]; !ok {
			return false
		}
	}
	return true
}

func (p *UintSet) Pack() []byte {
	buf := new(bytes.Buffer)
	tools.Dump(buf, uint32(len(p.vals)))
	for it, _ := range p.vals {
		tools.Dump(buf, it)
	}
	tools.Dump(buf, uint32(0))
	return buf.Bytes()
}

func (p *UintSet) Commit() []byte {
	buf := new(bytes.Buffer)
	tools.Dump(buf, uint32(len(p.adds)))
	for _, it := range p.adds {
		tools.Dump(buf, it)
	}
	p.adds = nil
	tools.Dump(buf, uint32(len(p.dels)))
	for _, it := range p.dels {
		tools.Dump(buf, it)
	}
	p.dels = nil
	return buf.Bytes()
}

func (p *UintSet) Merge(delta []byte) {
	r := bytes.NewReader(delta)
	ac := tools.Loadu32(r)
	for i := uint32(0); i < ac; i ++ {
		val := tools.Loadu32(r)
		p.vals[val] = true
	}
	dc := tools.Loadu32(r)
	for i := uint32(0); i < dc; i ++ {
		val := tools.Loadu32(r)
		delete(p.vals, val)
	}
}

func NewUintSet() *UintSet {
	return &UintSet{make(map[uint32]bool), nil, nil}
}

type UintSet struct {
	vals map[uint32]bool
	adds []uint32
	dels []uint32
}

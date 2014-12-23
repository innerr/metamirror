package structs

import (
	"fmt"
	"io"
	"defynetwork.com/tools"
)

func (p Clocks) Load(r io.Reader) {
	c := tools.Loadu16(r)
	for i := uint16(0); i < c; i++ {
		k := tools.Loadu64(r)
		v := tools.Loadu32(r)
		p[k] = v
	}
}

func (p Clocks) Dump(w io.Writer) {
	tools.Dump(w, uint16(len(p)))
	for k, v := range p {
		tools.Dump(w, k)
		tools.Dump(w, v)
	}
}

func (p Clocks) Absorb(x Clocks) {
	for k, v2 := range x {
		v1, ok := p[k]
		if !ok || v2 > v1 {
			p[k] = v2
		}
	}
}

func (p Clocks) Edit(hid uint64) {
	m := uint32(0)
	for _, v := range p {
		if v > m {
			m = v
		}
	}
	p[hid] = m + 1
}

func (p Clocks) After(x Clocks) bool {
	m := p.Max()
	c := len(p)
	if k, ok := x[m]; !ok {
		c -= 1
	} else {
		if k >= p[m] {
			return false
		}
	}
	if c != len(x) {
		return false
	}
	for k, v2 := range x {
		if k == m {
			continue
		}
		v1, ok := p[k]
		if !ok || v1 != v2 {
			return false
		}
	}
	return true
}

func (p Clocks) Copy() Clocks {
	x := make(Clocks)
	for k, v := range p {
		x[k] = v
	}
	return x
}

func (p Clocks) Sig(sep string) string {
	mk, mv := p.max()
	return fmt.Sprintf("%v" + sep + "%v" + sep + "%x", len(p), mv, mk)
}

func (p Clocks) Max() uint64 {
	mk, _ := p.max()
	return mk
}

func (p Clocks) Compare(x Clocks) int {
	l1 := len(p)
	l2 := len(x)

	if l1 == l2 {
		c := Equal
		for k, v1 := range p {
			v2, ok := x[k]
			if !ok {
				return Conflicted
			}
			if v1 == v2 {
				continue
			}
			if v1 < v2 {
				if c == Greater {
					return Conflicted
				} else {
					c = Smaller
				}
			}
			if v1 > v2 {
				if c == Smaller {
					return Conflicted
				} else {
					c = Greater
				}
			}
		}
		return c
	}

	if l1 < l2 {
		for k, v1 := range p {
			v2, ok := x[k]
			if !ok {
				return Conflicted
			}
			if v1 == v2 {
				continue
			}
			if v1 > v2 {
				return Conflicted
			}
		}
		return Smaller
	}

	if l1 > l2 {
		for k, v2 := range x {
			v1, ok := p[k]
			if !ok {
				return Conflicted
			}
			if v1 == v2 {
				continue
			}
			if v1 < v2 {
				return Conflicted
			}
		}
		return Greater
	}

	panic("unexpected")
}

func (p Clocks) max() (uint64, uint32) {
	var mk uint64
	var mv uint32
	for k, v := range p {
		if v > mv {
			mk = k
			mv = v
		}
	}
	return mk, mv
}

func NewClocks() Clocks {
	return make(Clocks)
}

type Clocks map[uint64]uint32

const (
	Equal = 0
	Smaller = -1
	Greater = 1
	Conflicted = -2
)

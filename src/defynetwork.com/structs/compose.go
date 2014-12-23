package structs

import (
	"bytes"
	"defynetwork.com/tools"
)

func (p *Compose) Pack(clocks Clocks) Blob {
	w := new(bytes.Buffer)
	modified := false

	data := p.def.Pack()
	if data == nil {
		tools.Dump(w, uint16(0))
	} else {
		modified = true
		tools.Dump(w, uint16(1))
		tools.Dumpd(w, data)
	}

	tools.Dump(w, uint16(len(p.named)))
	for name, it := range p.named {
		data := it.Pack()
		if data != nil {
			modified = true
		}
		tools.Dumps(w, name)
		tools.Dumpd(w, data)
	}

	if !modified {
		return Blob{nil, nil}
	}
	return Blob{clocks, w.Bytes()}
}

func (p *Compose) Commit(clocks Clocks) Blob {
	w := new(bytes.Buffer)
	modified := false

	data := p.def.Commit()
	if data == nil {
		tools.Dump(w, uint16(0))
	} else {
		modified = true
		tools.Dump(w, uint16(1))
		tools.Dumpd(w, data)
	}

	tools.Dump(w, uint16(len(p.named)))
	for name, it := range p.named {
		data := it.Commit()
		if data != nil {
			modified = true
		}
		tools.Dumps(w, name)
		tools.Dumpd(w, data)
	}

	if !modified {
		return Blob{nil, nil}
	}
	return Blob{clocks, w.Bytes()}
}

func (p *Compose) Merge(blob Blob) {
	r := bytes.NewReader(blob.Data)

	flag := tools.Loadu16(r)
	if flag == uint16(1) {
		data := tools.Loadd(r)
		p.def.Merge(data)
	} else if flag != 0 {
		panic("wrong compose flag")
	}

	count := tools.Loadu16(r)
	for i := uint16(0); i < count; i++ {
		name := tools.Loads(r)
		data := tools.Loadd(r)
		if named, ok := p.named[name]; ok && len(data) != 0 {
			named.Merge(data)
		}
	}
}

func (p *Compose) Data() IData {
	return p.def
}

func (p *Compose) AddNamedData(name string, data IData) {
	p.named[name] = data
}

func (p *Compose) GetNamedData(name string) IData {
	return p.named[name]
}

func NewCompose(data IData) *Compose {
	return &Compose{data, make(map[string]IData)}
}

type Compose struct {
	def IData
	named map[string]IData
}


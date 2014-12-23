package structs

import (
	"io"
	"defynetwork.com/tools"
)

func (p *Blob) Load(r io.Reader) {
	p.Vcs.Load(r)
	c := tools.Loadu16(r)
	if c == uint16(0) {
		return
	}
	p.Data = make([]byte, c)
	_, err := r.Read(p.Data)
	if err != nil {
		panic(err)
	}
}

func (p Blob) Dump(w io.Writer) {
	p.Vcs.Dump(w)
	tools.Dump(w, uint16(len(p.Data)))
	_, err := w.Write(p.Data)
	if err != nil {
		panic(err)
	}
}

func (p *Blob) IsNil() bool {
	return len(p.Vcs) == 0 && p.Data == nil
}

func NewBlob() Blob {
	return Blob{NewClocks(), nil}
}

type Blob struct {
	Vcs Clocks
	Data []byte
}

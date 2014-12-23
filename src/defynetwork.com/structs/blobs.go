package structs

import (
	"io"
	"defynetwork.com/tools"
)

func (p BlobMap) Load(r io.Reader) {
	c := tools.Loadu16(r)
	for i := uint16(0); i < c; i++ {
		k := tools.Loadu64(r)
		b := new(Blobs)
		b.Load(r)
		p[k] = *b
	}
}

func (p BlobMap) Dump(w io.Writer) {
	tools.Dump(w, uint16(len(p)))
	for k, v := range p {
		tools.Dump(w, k)
		v.Dump(w)
	}
}

type BlobMap map[uint64]Blobs

func (p *Delta) Load(r io.Reader) {
	v := (*Blobs)(p)
	v.Load(r)
}

func (p *Delta) Dump(w io.Writer) {
	v := (*Blobs)(p)
	v.Dump(w)
}

type Delta Blobs

func (p *Blobs) Load(r io.Reader) {
	c := tools.Loadu16(r)
	for i := uint16(0); i < c; i++ {
		blob := NewBlob()
		blob.Load(r)
		*p = append(*p, blob)
	}
}

func (p *Blobs) Dump(w io.Writer) {
	tools.Dump(w, uint16(len(*p)))
	for _, blob := range *p {
		blob.Dump(w)
	}
}

type Blobs []Blob

package structs

import (
	"sort"
)

func (p *Box) Delta(clocks Clocks) Delta {
	delta := Delta{}
	for hid, blobs := range p.blobs {
		finder := _NewBlobsFinder(hid, clocks, blobs)
		delta = append(delta, finder.Delta()...)
	}
	return delta
}

func (p *Box) Merge(delta Delta) {
	for _, blob := range delta {
		p.Add(blob, blob.Vcs.Max())
	}
}

func (p *Box) Edit(hid uint64) Clocks {
	clocks := p.max.Copy()
	clocks.Edit(hid)
	return clocks
}

func (p *Box) Add(blob Blob, hid uint64) {
	if blobs, ok := p.blobs[hid]; !ok {
		p.blobs[hid] = Blobs{blob}
	} else {
		p.blobs[hid] = append(blobs, blob)
	}
	p.max.Absorb(blob.Vcs)
}

func (p *Box) Max() Clocks {
	return p.max.Copy()
}

func NewBox() *Box {
	return &Box{make(BlobMap), make(Clocks)}
}

type Box struct {
	blobs BlobMap
	max Clocks
}

type _BlobsFinder struct {
	hid uint64
	clocks Clocks
	blobs Blobs
	ver uint32
}

func _NewBlobsFinder(hid uint64, clocks Clocks, blobs Blobs) _BlobsFinder {
	return _BlobsFinder{hid, clocks, blobs, clocks[hid]}
}

func (p _BlobsFinder) Delta() Blobs {
	i := sort.Search(len(p.blobs), p.finder)
	if i >= len(p.blobs) {
		return nil
	}
	return p.blobs[i:]
}

func (p _BlobsFinder) finder(i int) bool {
	blob := p.blobs[i]
	if blob.Vcs[p.hid] >= p.ver {
		c := blob.Vcs.Compare(p.clocks)
		if c == Greater || c == Conflicted {
			return true
		}
	}
	return false
}

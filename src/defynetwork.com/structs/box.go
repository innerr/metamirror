package structs

import "sort"

func (p *Box) Delta(clocks Clocks) Delta {
	delta := Delta{}
	for hid, blobs := range p.blobs {
		a := clocks[hid]
		i := sort.Search(len(blobs), func(i int) bool {
			b := blobs[i]
			if b.Vcs[hid] >= a {
				c := b.Vcs.Compare(clocks)
				if c == Greater || c == Conflicted {
					return true
				}
			}
			return false
		})
		if i >= len(blobs) {
			continue
		}
		delta = append(delta, blobs[i:]...)
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

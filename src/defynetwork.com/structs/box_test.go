package structs

import (
	"testing"
	"defynetwork.com/tools"
)

func TestBox(t *testing.T) {
	merge := func(c *Box, d Blob) {
		c.Merge(Delta{d})
	}

	delta := func(blobs ...Blob) Delta {
		m := NewBox()
		for _, blob := range blobs {
			m.Add(blob, blob.Vcs.Max())
		}
		return m.Delta(Clocks{})
	}

	check := func(d1, d2 Delta) {
		tom := func(d Delta) BlobMap {
			bm := BlobMap{}
			for _, it := range d{
				hid := it.Vcs.Max()
				bs, ok := bm[hid]
				if ok {
					bm[hid] = append(bs, it)
				} else {
					bm[hid] = Blobs{it}
				}
			}
			return bm
		}
		tools.Check(tom(d1), tom(d2))
	}

	id1 := uint64(111)
	id2 := uint64(222)

	d1 := []byte("d1")
	d2 := []byte("d2")
	d3 := []byte("d3")

	c := NewBox()
	v := c.Edit(id1)
	tools.Check(v, Clocks{id1: 1})
	c.Add(Blob{v, d1}, id1)

	v = c.Edit(id2)
	tools.Check(v, Clocks{id1: 1, id2: 2})
	merge(c, Blob{v, d2})

	v = c.Edit(id1)
	tools.Check(v, Clocks{id1: 3, id2: 2})
	merge(c, Blob{v, d3})

	b1 := Blob{Clocks{id1: 1}, d1}
	b2 := Blob{Clocks{id1: 1, id2: 2}, d2}
	b3 := Blob{Clocks{id1: 3, id2: 2}, d3}

	d := c.Delta(Clocks{id1: 3, id2: 2})
	check(d, delta())

	d = c.Delta(Clocks{id1: 2, id2: 2})
	check(d, delta(b3))

	d = c.Delta(Clocks{id1: 1, id2: 2})
	check(d, delta(b3))

	d = c.Delta(Clocks{id1: 2, id2: 1})
	check(d, delta(b2, b3))

	d = c.Delta(Clocks{id1: 1, id2: 1})
	check(d, delta(b2, b3))

	d = c.Delta(Clocks{id1: 1})
	check(d, delta(b2, b3))

	d = c.Delta(Clocks{id2: 1})
	check(d, delta(b1, b2, b3))

	d = c.Delta(Clocks{})
	check(d, delta(b1, b2, b3))
}


package structs

import "testing"

func TestCompose(t *testing.T) {
	d := NewUintSet()
	c := NewCompose(d)
	d.Set(6)
	d.Set(9)
	b := c.Commit(Clocks{})

	x := NewCompose(NewUintSet())
	x.Merge(b)
	u := x.Data().(*UintSet)
	if !u.Has(6) || !u.Has(9) {
		t.Fatal("wrong")
	}
}

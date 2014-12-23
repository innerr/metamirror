package structs

import (
	"bytes"
	"testing"
	"defynetwork.com/tools"
)

func TestClocks(t *testing.T) {
	id1 := uint64(111)
	id2 := uint64(222)

	v1 := Clocks{id1: 1}
	v1.Edit(id1)
	tools.Check(v1, Clocks{id1: 2})
	v1.Edit(id1)
	tools.Check(v1, Clocks{id1: 3})

	v2 := v1.Copy()
	v2.Edit(id2)
	tools.Check(v2, Clocks{id1: 3, id2: 4})

	v2 = Clocks{id1: 1, id2: 2}
	v2.Edit(id1)
	tools.Check(v2, Clocks{id1: 3, id2: 2})

	v2 = v1.Copy()
	tools.Check(v2, Clocks{id1: 3})

	v1 = Clocks{id1: 2, id2: 1}
	v2 = Clocks{id1: 2, id2: 4}
	v1.Absorb(v2)
	tools.Check(v2, Clocks{id1: 2, id2: 4})
	v1 = Clocks{id1: 2, id2: 4}
	v2 = Clocks{id1: 2, id2: 1}
	v1.Absorb(v2)
	tools.Check(v1, Clocks{id1: 2, id2: 4})
	v1 = Clocks{id1: 2, id2: 1}
	v2 = Clocks{id2: 4}
	v1.Absorb(v2)
	tools.Check(v1, Clocks{id1: 2, id2: 4})
	v1 = Clocks{id1: 2}
	v2 = Clocks{id2: 4}
	v1.Absorb(v2)
	tools.Check(v1, Clocks{id1: 2, id2: 4})

	v1 = Clocks{id1: 1}
	v2 = Clocks{id1: 1}
	tools.Check(v1.Compare(v2), Equal)
	tools.Check(v2.Compare(v1), Equal)
	v1.Edit(id1)
	tools.Check(v1.Compare(v2), Greater)
	tools.Check(v2.Compare(v1), Smaller)
	v2.Edit(id2)
	tools.Check(v1.Compare(v2), Conflicted)
	tools.Check(v2.Compare(v1), Conflicted)

	v1 = Clocks{id1: 2, id2: 1}
	v2 = Clocks{id1: 2, id2: 1}
	tools.Check(v1.Compare(v2), Equal)
	tools.Check(v2.Compare(v1), Equal)
	v1 = Clocks{id1: 2, id2: 2}
	v2 = Clocks{id1: 2, id2: 1}
	tools.Check(v1.Compare(v2), Greater)
	tools.Check(v2.Compare(v1), Smaller)
	v1 = Clocks{id1: 2, id2: 2}
	v2 = Clocks{id1: 2}
	tools.Check(v1.Compare(v2), Greater)
	tools.Check(v2.Compare(v1), Smaller)
	v1 = Clocks{id1: 2, id2: 2}
	v2 = Clocks{id1: 3}
	tools.Check(v1.Compare(v2), Conflicted)
	tools.Check(v2.Compare(v1), Conflicted)

	v1 = Clocks{id1: 2, id2: 1}
	v2 = Clocks{id1: 1, id2: 2}
	tools.Check(v1.Max(), id1)
	tools.Check(v2.Max(), id2)

	v1 = Clocks{id1: 1}
	v2 = Clocks{id1: 1}
	v3 := v1.Copy()
	v1.Edit(id1)
	tools.Check(v1.After(v3), true)
	tools.Check(v2.After(v3), false)
	v4 := v2.Copy()
	v2.Edit(id2)
	tools.Check(v1.After(v4), true)
	tools.Check(v2.After(v4), true)

	v1 = Clocks{id1: 3, id2: 1}
	buf := new(bytes.Buffer)
	v1.Dump(buf)
	v2 = Clocks{}
	v2.Load(bytes.NewReader(buf.Bytes()))
	tools.Check(v1, v2)
}


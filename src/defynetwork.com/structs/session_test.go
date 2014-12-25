package structs

import (
	"testing"
	"defynetwork.com/tools"
)

func TestSession(t *testing.T) {
	testSession(
		&CoreFlags{In: true, Out: true},
		&CoreFlags{In: true, Out: true},
		4,
		9,
		NewUintSet().Set(4).Set(9),
		NewUintSet().Set(4).Set(9),
		t)

	testSession(
		&CoreFlags{In: false, Out: false},
		&CoreFlags{In: false, Out: false},
		4,
		9,
		NewUintSet().Set(4),
		NewUintSet().Set(9),
		t)

	testSession(
		&CoreFlags{In: true, Out: false},
		&CoreFlags{In: false, Out: false},
		4,
		9,
		NewUintSet().Set(4),
		NewUintSet().Set(9),
		t)

	testSession(
		&CoreFlags{In: false, Out: true},
		&CoreFlags{In: false, Out: false},
		4,
		9,
		NewUintSet().Set(4),
		NewUintSet().Set(9),
		t)

	testSession(
		&CoreFlags{In: false, Out: false},
		&CoreFlags{In: true, Out: false},
		4,
		9,
		NewUintSet().Set(4),
		NewUintSet().Set(9),
		t)

	testSession(
		&CoreFlags{In: false, Out: false},
		&CoreFlags{In: false, Out: true},
		4,
		9,
		NewUintSet().Set(4),
		NewUintSet().Set(9),
		t)

	testSession(
		&CoreFlags{In: true, Out: true},
		&CoreFlags{In: false, Out: false},
		4,
		9,
		NewUintSet().Set(4),
		NewUintSet().Set(9),
		t)

	testSession(
		&CoreFlags{In: true, Out: false},
		&CoreFlags{In: true, Out: false},
		4,
		9,
		NewUintSet().Set(4),
		NewUintSet().Set(9),
		t)

	testSession(
		&CoreFlags{In: true, Out: false},
		&CoreFlags{In: false, Out: true},
		4,
		9,
		NewUintSet().Set(4).Set(9),
		NewUintSet().Set(9),
		t)

	testSession(
		&CoreFlags{In: false, Out: true},
		&CoreFlags{In: true, Out: false},
		4,
		9,
		NewUintSet().Set(4),
		NewUintSet().Set(9).Set(4),
		t)

	testSession(
		&CoreFlags{In: false, Out: false},
		&CoreFlags{In: true, Out: true},
		4,
		9,
		NewUintSet().Set(4),
		NewUintSet().Set(9),
		t)

	testSession(
		&CoreFlags{In: true, Out: true},
		&CoreFlags{In: true, Out: false},
		4,
		9,
		NewUintSet().Set(4),
		NewUintSet().Set(9).Set(4),
		t)

	testSession(
		&CoreFlags{In: false, Out: true},
		&CoreFlags{In: true, Out: true},
		4,
		9,
		NewUintSet().Set(4),
		NewUintSet().Set(9).Set(4),
		t)

	testSession(
		&CoreFlags{In: true, Out: false},
		&CoreFlags{In: true, Out: true},
		4,
		9,
		NewUintSet().Set(4).Set(9),
		NewUintSet().Set(9),
		t)

	testSession(
		&CoreFlags{In: true, Out: true},
		&CoreFlags{In: false, Out: true},
		4,
		9,
		NewUintSet().Set(4).Set(9),
		NewUintSet().Set(9),
		t)
}

func testSession(f1, f2 *CoreFlags, k1, k2 uint32, r1, r2 *UintSet, t *testing.T) {
	l1 := tools.NewLog("", false, tools.LogLvlDetail).Mod("sync").Name("l1")
	l2 := tools.NewLog("", false, tools.LogLvlDetail).Mod("sync").Name("l2")
	persist := &DumbPersist{}
	c := NewBiChannel()

	df := make(chan bool, 2)
	done := func() {
		df <-true
	}

	u1 := NewUintSet()
	c1 := NewCore(u1, persist, f1, l1)
	s1 := NewSession(c1, c.A, nil, done, l1)
	u1.Set(k1)
	c1.Commit(11)

	u2 := NewUintSet()
	c2 := NewCore(u2, persist, f2, l2)
	s2 := NewSession(c2, c.B, nil, done, l2)
	u2.Set(k2)
	c2.Commit(22)

	s2.Sync()
	<-df
	<-df
	if !u1.Equal(r1) {
		t.Fatal("wrong1")
	}

	s1.Sync()
	<-df
	<-df
	if !u2.Equal(r2) {
		t.Fatal("wrong2")
	}
}

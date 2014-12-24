package structs

import (
	"os"
	"testing"
	"defynetwork.com/tools"
)

func TestCoreMerge(t *testing.T) {
	log := tools.NewLog("", false, tools.LogLvlNone)
	persist := &DumbPersist{}

	s1 := NewUintSet()
	c1 := NewCore(s1, persist, nil, log)
	s1.Set(4)
	s1.Set(10)
	d1 := c1.Commit(69)
	d2, _ := c1.Pack()

	s2 := NewUintSet()
	c2 := NewCore(s2, persist, nil, log)
	c2.Merge(d1)
	if !s1.Equal(s2) {
		t.Fatal("a != b")
	}

	s3 := NewUintSet()
	c3 := NewCore(s3, persist, nil, log)
	c3.Merge(d2)
	if !s1.Equal(s3) {
		t.Fatal("a != b")
	}
}

func TestCorePersist(t *testing.T) {
	path := "core-persist.test"
	defer os.Remove(path)

	log := tools.NewLog("", false, tools.LogLvlNone)

	s1 := NewUintSet()
	f1, err := os.OpenFile(path, os.O_RDWR | os.O_CREATE | os.O_SYNC, 0640)
	if err != nil {
		t.Fatal(err)
	}
	p1 := NewPersist(f1)
	c := NewCore(s1, p1, nil, log)
	s1.Set(4)
	s1.Set(10)
	c.Commit(69)

	s2 := NewUintSet()
	f2, err := os.OpenFile(path, os.O_RDWR | os.O_CREATE | os.O_SYNC, 0640)
	if err != nil {
		t.Fatal(err)
	}
	p2 := NewPersist(f2)
	c = NewCore(s2, p2, nil, log)
	if !s1.Equal(s2) {
		t.Fatal("a != b")
	}
}

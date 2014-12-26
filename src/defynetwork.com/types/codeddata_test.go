package structs

import (
	"testing"
)

func TestCodedData(t *testing.T) {
	key := []byte("testtesttesttest")

	s1 := NewUintSet()
	c1 := NewCodedData(s1, true, key)
	s1.Set(4)
	s1.Set(10)
	d1 := c1.Commit()
	d2 := c1.Pack()

	s2 := NewUintSet()
	c2 := NewCodedData(s2, true, key)
	c2.Merge(d1)
	if !s1.Equal(s2) {
		t.Fatal("a != b")
	}

	s3 := NewUintSet()
	c3 := NewCodedData(s3, true, key)
	c3.Merge(d2)
	if !s1.Equal(s3) {
		t.Fatal("a != b")
	}
}

func TestCodedDataBenchmark(t *testing.T) {
	// TODO: coded data benchmark
}


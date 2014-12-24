package structs

import "testing"

func TestUintSet(t *testing.T) {
	s1 := NewUintSet()
	s1.Set(4)
	s1.Set(10)
	d1 := s1.Commit()
	d2 := s1.Pack()

	s2 := NewUintSet()
	s2.Merge(d1)
	if !s1.Equal(s2) {
		t.Fatal("a != b")
	}

	s3 := NewUintSet()
	s3.Merge(d2)
	if !s1.Equal(s3) {
		t.Fatal("a != b")
	}
}

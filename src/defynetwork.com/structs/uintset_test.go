package structs

import "testing"

func TestUintSet(t *testing.T) {
	s1 := NewUintSet()
	s1.Set(4)
	s1.Set(10)
	d := s1.Commit()

	s2 := NewUintSet()
	s2.Merge(d)
	if !s1.Equal(s2) {
		t.Fatal("a != b")
	}
}

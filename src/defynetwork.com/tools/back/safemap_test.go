package tools

import "testing"

func TestSaftMap(t *testing.T) {
	m := NewSafeMap()
	m.Set("abc", "a")
	m.Set(4, "b")
	m.Set(5, 6)

	v, ok := m.Get("abc")
	Check(ok, true)
	Check(v.(string), "a")
	v, ok = m.Get(4)
	Check(ok, true)
	Check(v.(string), "b")
	v, ok = m.Get(5)
	Check(ok, true)
	Check(v.(int), 6)
	v, ok = m.Get(6)
	Check(ok, false)
}

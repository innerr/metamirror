package structs

import (
	"bytes"
	"io"
	"io/ioutil"
	"testing"
)

func TestRpc(t *testing.T) {
	c := NewBiChannel()

	var a, b string
	var an, bn int

	ra := NewRpc(c.A)
	ra.Func("ba").Receive(func(r io.Reader, n uint32) {
		d, _ := ioutil.ReadAll(r)
		a = string(d)
		an = len(d)
	})

	rb := NewRpc(c.B)
	rb.Func("ab").Receive(func(r io.Reader, n uint32) {
		d, _ := ioutil.ReadAll(r)
		b = string(d)
		bn = len(d)
	})

	as := []byte("send from a")
	ra.Func("ab").Send(bytes.NewReader(as), uint32(len(as)))
	bs := []byte("send from b, test")
	rb.Func("ba").Send(bytes.NewReader(bs), uint32(len(bs)))
	c.Close()

	if b != string(as) || bn != len(as) {
		t.Fatal(b, "!=", string(as))
	}
	if a != string(bs) || an != len(bs) {
		t.Fatal(a, "!=", string(bs))
	}
}

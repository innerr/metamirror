package structs

import (
	"bytes"
	"io"
	"io/ioutil"
	"testing"
)

func TestBiChannel(t *testing.T) {
	c := NewBiChannel()

	var a, b string
	var an, bn int
	c.A.Receive(func(r io.Reader, n uint32) {
		d, _ := ioutil.ReadAll(r)
		a = string(d)
		an = len(d)
	})
	c.B.Receive(func(r io.Reader, n uint32) {
		d, _ := ioutil.ReadAll(r)
		b = string(d)
		bn = len(d)
	})

	as := []byte("send from a")
	c.A.Send(bytes.NewReader(as), uint32(len(as)))
	bs := []byte("send from b, test")
	c.B.Send(bytes.NewReader(bs), uint32(len(bs)))
	c.Close()

	if a != string(bs) || an != len(bs) {
		t.Fatal(a, "!=", bs)
	}
	if b != string(as) || bn != len(as) {
		t.Fatal(b, "!=", as)
	}
}

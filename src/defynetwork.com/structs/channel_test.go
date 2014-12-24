package structs

import (
	"bytes"
	"io"
	"io/ioutil"
	"testing"
)

func TestAsynChannel(t *testing.T) {
	var res []string
	mc := NewMemChannel(nil, func(r io.Reader, n uint32){
		d, _ := ioutil.ReadAll(r)
		if n != uint32(len(d)) {
			t.Fatal("wrong")
		}
		res = append(res, string(d))
	})

	mc.Send(bytes.NewReader([]byte("1")), 1)
	if len(res) != 1 || res[0] != "1" {
		t.Fatal("wrong")
	}

	ac := NewAsynChannel(mc, 0)
	mc.Send(bytes.NewReader([]byte("2")), 1)
	ac.Close()
	if len(res) != 2 || res[1] != "2" {
		t.Fatal("wrong")
	}
}

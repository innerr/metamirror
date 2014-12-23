package structs

import (
	"bytes"
	"testing"
	"defynetwork.com/tools"
)

func TestBlob(t *testing.T) {
	b1 := Blob {
		Clocks{11: 2, 22: 4},
		[]byte("test b1"),
	}
	buf := new(bytes.Buffer)
	b1.Dump(buf)

	b2 := NewBlob()
	b2.Load(bytes.NewReader(buf.Bytes()))
	tools.Check(b1, b2)

	buf.Reset()
	b3 := Blob {
		Clocks{14: 9, 28: 3},
		[]byte("test b3"),
	}
	bs1 := Blobs{b1, b3, b2}
	bs1.Dump(buf)
	bs2 := Blobs{}
	bs2.Load(bytes.NewReader(buf.Bytes()))
	tools.Check(bs1, bs2)
}

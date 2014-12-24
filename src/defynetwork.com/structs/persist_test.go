package structs

import (
	"bytes"
	"testing"
	"defynetwork.com/tools"
)

func TestMemPersist(t *testing.T) {
	p := NewPersist(new(bytes.Buffer))
	a := Blobs {
		Blob {
			Clocks{11: 2, 22: 4},
			[]byte("test b1"),
		},
		Blob {
			Clocks{33: 2, 22: 3},
			[]byte("test b2"),
		},
		Blob {
			Clocks{44: 1, 1: 6},
			[]byte("test b3"),
		},
	}
	for _, it := range a {
		p.Dump(it)
	}
	b := Blobs{}
	p.Load(func(blob Blob) {
		b = append(b, blob)
	})
	tools.Check(a, b)
}

func TestPersistBenchmark(t *testing.T) {
	// TODO: persist benchmark
}

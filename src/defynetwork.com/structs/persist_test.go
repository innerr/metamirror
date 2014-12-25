package structs

import (
	"os"
	"testing"
	"defynetwork.com/tools"
)

func TestPersist(t *testing.T) {
	path := "persist.test"
	defer os.Remove(path)

	log := tools.NewLog("", false, tools.LogLvlNone)

	p1 := NewPersistWithPath(path, log)
	defer p1.Close()
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
		p1.Dump(it)
	}

	b := Blobs{}
	p2 := NewPersistWithPath(path, log)
	defer p2.Close()
	p2.Load(func(blob Blob) {
		b = append(b, blob)
	})
	tools.Check(a, b)
}

func TestPersistBenchmark(t *testing.T) {
	// TODO: persist benchmark
}

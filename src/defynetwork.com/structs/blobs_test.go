package structs

import (
	"bytes"
	"testing"
	"defynetwork.com/tools"
)

func TestBlobMap(t *testing.T) {
	b1 := BlobMap {
		111: Blobs {
			Blob {
				Clocks{11: 2, 22: 4},
				[]byte("test b1"),
			},
		},
		222 : Blobs {
			Blob {
				Clocks{33: 2, 22: 4},
				[]byte("test b2"),
			},
			Blob {
				Clocks{33: 2, 4: 4},
				[]byte("test b3"),
			},
		},
	}
	buf := new(bytes.Buffer)
	b1.Dump(buf)

	b2 := BlobMap{}
	b2.Load(bytes.NewReader(buf.Bytes()))
	tools.Check(b1, b2)
}

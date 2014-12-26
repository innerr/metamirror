package tools

import (
	"bytes"
	"testing"
)

func TestIoUtil(t *testing.T) {
	w := new(bytes.Buffer)
	Dumpb(w, true)
	Dumpb(w, false)
	Dumps(w, "test")
	Dumpd(w, []byte("test"))
	Dump(w, uint16(1))
	Dump(w, uint32(2))
	Dump(w, uint64(3))
	Dump(w, int16(1))
	Dump(w, int32(2))
	Dump(w, int64(3))

	r := bytes.NewReader(w.Bytes())
	Check(Loadb(r), true)
	Check(Loadb(r), false)
	Check(Loads(r), "test")
	Check(Loadd(r), []byte("test"))
	Check(Loadu16(r), uint16(1))
	Check(Loadu32(r), uint32(2))
	Check(Loadu64(r), uint64(3))
	Check(Loadn16(r), int16(1))
	Check(Loadn32(r), int32(2))
	Check(Loadn64(r), int64(3))
}

package tools

import (
	"bytes"
	"encoding/binary"
	"io"
)

func Dump(w io.Writer, v interface{}) {
	err := binary.Write(w, binary.LittleEndian, v)
	if err != nil {
		panic(err)
	}
}

func Load(r io.Reader, v interface{}) {
	err := binary.Read(r, binary.LittleEndian, v)
	if err != nil {
		panic(err)
	}
}

func Loadu16(r io.Reader) uint16 {
	n := uint16(0)
	Load(r, &n)
	return n
}

func Loadu32(r io.Reader) uint32 {
	n := uint32(0)
	Load(r, &n)
	return n
}

func Loadu64(r io.Reader) uint64 {
	n := uint64(0)
	Load(r, &n)
	return n
}

func Loadn16(r io.Reader) int16 {
	n := int16(0)
	Load(r, &n)
	return n
}

func Loadn32(r io.Reader) int32 {
	n := int32(0)
	Load(r, &n)
	return n
}

func Loadn64(r io.Reader) int64 {
	n := int64(0)
	Load(r, &n)
	return n
}

func Dumpd(w io.Writer, d []byte) {
	Dump(w, uint32(len(d)))
	_, err := w.Write(d)
	if err != nil {
		panic(err)
	}
}

func Loadd(r io.Reader) []byte {
	c := uint32(0)
	Load(r, &c)
	d := make([]byte, c)
	_, err := r.Read(d)
	if err != nil {
		panic(err)
	}
	return d
}

func Dumpb(w io.Writer, v bool) {
	if v {
		Dump(w, uint16(1))
	} else {
		Dump(w, uint16(0))
	}
}

func Loadb(r io.Reader) bool {
	n := uint16(0)
	Load(r, &n)
	return n == 1
}

func Dumps(w io.Writer, s string) {
	d := []byte(s)
	Dump(w, uint16(len(d)))
	_, err := w.Write(d)
	if err != nil {
		panic(err)
	}
}

func Loads(r io.Reader) string {
	c := uint16(0)
	Load(r, &c)
	d := make([]byte, c)
	_, err := r.Read(d)
	if err != nil {
		panic(err)
	}
	return string(d)
}

func Packss(strs []string) []byte {
	buf := new(bytes.Buffer)
	Dump(buf, uint16(len(strs)))
	for _, str := range strs {
		Dumps(buf, str)
	}
	return buf.Bytes()
}

func Unpackss(r io.Reader) []string {
	c := Loadu16(r)
	strs := make([]string, c)
	for i := uint16(0); i < c; i++ {
		strs[i] = Loads(r)
	}
	return strs
}

func Packs(s string) []byte {
	w := new(bytes.Buffer)
	Dumps(w, s)
	return w.Bytes()
}

func Packd(data []byte) []byte {
	buf := new(bytes.Buffer)
	Dumpd(buf, data)
	return buf.Bytes()
}

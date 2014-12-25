package structs

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"hash/crc32"
	"os"
	"defynetwork.com/tools"
)

type PersistIOError struct {
	error
}

func (p *Persist) Load(fun func(blob Blob)) {
	if p.file == nil {
		return
	}

	p.log.Debug("loading")

	var sum int64
	var size, crc uint32

	for ; sum < p.fsize; sum += int64(size + 4 + 4) {
		defer func() {
			if e := recover(); e != nil {
				if err, ok := e.(error); ok && err == io.EOF {
					return
				}
				panic(e)
			}
		}()

		err := binary.Read(p.file, binary.LittleEndian, &size)
		if err != nil {
			panic(&PersistIOError{err})
		}
		err = binary.Read(p.file, binary.LittleEndian, &crc)
		if err != nil {
			panic(&PersistIOError{err})
		}

		data := make([]byte, size)
		_, err = io.ReadFull(p.file, data)
		if err != nil {
			panic(&PersistIOError{err})
		}
		digest := crc32.ChecksumIEEE(data)
		if digest != crc {
			panic(&PersistIOError{fmt.Errorf("crc32 unmatched: %v != %v", crc, digest)})
		}

		blob := NewBlob()
		blob.Load(bytes.NewReader(data))
		fun(blob)
	}

	if sum != p.fsize {
		panic(&PersistIOError{fmt.Errorf("tail data left: %v - %v = %v", p.fsize, sum, p.fsize - sum)})
	}
	p.fsize = -1

	p.log.Debug("loaded")
}

func (p *Persist) Dump(blob Blob) {
	if p.file == nil {
		return
	}

	buf1 := new(bytes.Buffer)
	blob.Dump(buf1)
	data := buf1.Bytes()

	crc := crc32.ChecksumIEEE(data)
	buf2 := new(bytes.Buffer)
	tools.Dump(buf2, uint32(len(data)))
	tools.Dump(buf2, crc)
	header := buf2.Bytes()

	_, err := p.file.Write(header)
	if err != nil {
		panic(err)
	}
	_, err = p.file.Write(data)
	if err != nil {
		panic(err)
	}
}

func (p *Persist) Mark() int64 {
	mark, err := p.file.Seek(0, os.SEEK_CUR)
	if err != nil {
		panic(err)
	}
	return mark
}

func (p *Persist) Rollback(mark int64) {
	err := p.file.Truncate(mark)
	if err != nil {
		panic(err)
	}
}

func (p *Persist) Close() {
	err := p.file.Close()
	if err != nil {
		panic(err)
	}
}

func NewPersistWithPath(path string, log *tools.Log) *Persist {
	file, err := os.OpenFile(path, os.O_RDWR | os.O_CREATE, 0640)
	if err != nil {
		panic(err)
	}
	info, err := os.Stat(path)
	if err != nil {
		panic(err)
	}
	return NewPersist(file, info.Size(), log)
}

func NewPersist(file *os.File, fsize int64, log *tools.Log) *Persist {
	return &Persist{file, fsize, log}
}

type Persist struct {
	file *os.File
	fsize int64
	log *tools.Log
}

func (p *DumbPersist) Load(fun func(blob Blob)) {
}
func (p *DumbPersist) Dump(blob Blob) {
}
func (p *DumbPersist) Mark() int64 {
	return -1
}
func (p *DumbPersist) Rollback(int64) {
}
func (p *DumbPersist) Close() {
}
type DumbPersist struct {
}

type IPersist interface {
	Load(func(Blob))
	Dump(Blob)
	Mark() int64
	Rollback(int64)
	Close()
}

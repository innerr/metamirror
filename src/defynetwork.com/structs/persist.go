package structs

import "io"

func (p *DumbPersist) Load(fun func(blob Blob)) {
}
func (p *DumbPersist) Dump(blob Blob) {
}
type DumbPersist struct {
}

func (p *Persist) Load(fun func(blob Blob)) {
	if p.file == nil {
		return
	}
	for {
		defer func() {
			if e := recover(); e != nil {
				if err, ok := e.(error); ok && err == io.EOF {
					return
				}
				panic(e)
			}
		}()
		blob := NewBlob()
		blob.Load(p.file)
		fun(blob)
	}
}

func (p *Persist) Dump(blob Blob) {
	if p.file == nil {
		return
	}
	blob.Dump(p.file)
}

func NewPersist(file File) *Persist {
	return &Persist{file}
}

type Persist struct {
	file File
}

type File interface {
	io.Reader
	io.Writer
}

type IPersist interface {
	Load(func(Blob))
	Dump(Blob)
}

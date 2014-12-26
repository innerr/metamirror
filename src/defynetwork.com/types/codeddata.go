package structs

import (
	"bytes"
	"defynetwork.com/tools"
)

func (p *CodedData) Pack() []byte {
	data := p.origin.Pack()
	return tools.Encode(bytes.NewReader(data), uint32(len(data)), p.compress, p.keyword)
}

func (p *CodedData) Commit() []byte {
	data := p.origin.Commit()
	return tools.Encode(bytes.NewReader(data), uint32(len(data)), p.compress, p.keyword)
}

func (p *CodedData) Merge(data []byte) {
	decoded := tools.Decode(data, p.keyword)
	p.origin.Merge(decoded)
}

func NewCodedData(origin IData, compress bool, keyword []byte) * CodedData {
	return &CodedData{origin, compress, keyword}
}

type CodedData struct {
	origin IData
	compress bool
	keyword []byte
}

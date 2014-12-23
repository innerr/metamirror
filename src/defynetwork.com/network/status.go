package network

import (
	"encoding/hex"
	"sync"
	"time"
	"defynetwork.com/tools"
)

func (p *TransOuts) Count() int {
	p.lock.Lock()
	defer p.lock.Unlock()
	return len(p.outs)
}

func (p *TransOuts) Add(hash []byte) {
	p.lock.Lock()
	defer p.lock.Unlock()
	esha1 := hex.EncodeToString(hash)
	p.outs[esha1] = &TransOut{hash, time.Now().UnixNano()}
}

func (p *TransOuts) Done(hash []byte) {
	p.lock.Lock()
	defer p.lock.Unlock()
	esha1 := hex.EncodeToString(hash)
	delete(p.outs, esha1)
}

func NewTransOuts(log *tools.Log) *TransOuts {
	return &TransOuts{log, make(map[string]*TransOut), sync.Mutex{}}
}

type TransOuts struct {
	log *tools.Log
	outs map[string]*TransOut
	lock sync.Mutex
}

type TransOut struct {
	Hash []byte
	Atime int64
}

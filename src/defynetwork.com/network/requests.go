package network

import (
	"bytes"
	"io"
	"sync"
	"time"
	"defynetwork.com/tools"
)

func (p *Requests) Count() int {
	p.lock.Lock()
	defer p.lock.Unlock()
	return len(p.reqs)
}

func (p *Requests) Walk(fun func(req *Request)) {
	p.lock.Lock()
	defer p.lock.Unlock()
	now := time.Now().UnixNano()
	for _, req := range p.reqs {
		if time.Duration(now - req.Rtime) > p.timeout {
			fun(req)
		}
	}
}

func (p *Requests) Requesting(key string) {
	p.lock.Lock()
	defer p.lock.Unlock()
	req, ok := p.reqs[key]
	if ok {
		req.Rtime = time.Now().UnixNano()
	}
}

func (p *Requests) Add(key string, hash []byte, fun RequestDoneFunc) {
	p.lock.Lock()
	defer p.lock.Unlock()
	req, ok := p.reqs[key]
	if ok {
		req.Funs = append(req.Funs, fun)
		return
	}
	req = &Request{hash, []RequestDoneFunc{fun}, nil, time.Now().UnixNano(), -1, 0}
	p.reqs[key] = req
}

func (p *Requests) Done(key string, r io.Reader, n uint32) (total int, done int) {
	p.lock.Lock()
	defer p.lock.Unlock()
	req, ok := p.reqs[key]
	if !ok {
		return 0, 0
	}
	if len(req.Funs) == 1 {
		req.Funs[0](r, n)
	} else {
		buf := new(bytes.Buffer)
		read, err := io.CopyN(buf, r, int64(n))
		if err != nil {
			tools.Eat(p.log, r, n - uint32(read))
			req.Failed += 1
			return len(req.Funs), 0
		}
		data := buf.Bytes()
		for _, fun := range req.Funs {
			fun(bytes.NewReader(data), n)
		}
	}
	delete(p.reqs, key)
	return len(req.Funs), len(req.Funs)
}

// TODO: use more chs to get data
func (p *Requests) Found(key string, chs *Channels) int {
	p.lock.Lock()
	defer p.lock.Unlock()
	req, ok := p.reqs[key]
	if !ok {
		return 0
	}
	req.Found = append(req.Found, chs)
	return len(req.Found)
}

func NewRequests(log *tools.Log, timeout int) *Requests {
	return &Requests{log, make(map[string]*Request), time.Duration(timeout) * time.Millisecond, sync.Mutex{}}
}

type Requests struct {
	log *tools.Log
	reqs map[string]*Request
	timeout time.Duration
	lock sync.Mutex
}

type Request struct {
	Hash []byte
	Funs []RequestDoneFunc
	Found []*Channels
	Atime int64
	Rtime int64
	Failed int
}

// TODO: what if failed?
type RequestDoneFunc func(r io.Reader, n uint32)

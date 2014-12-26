package tools

import "sync"

func (p *SafeMap) Size() int {
	p.lock.Lock()
	defer p.lock.Unlock()
	return len(p.vals)
}

func (p *SafeMap) Del(k interface{}) {
	delete(p.vals, k)
}

func (p *SafeMap) Walk(fun func(k, v interface{})) {
	p.lock.Lock()
	defer p.lock.Unlock()
	for k, v := range p.vals {
		fun(k, v)
	}
}

func (p *SafeMap) Has(k interface{}) bool {
	p.lock.Lock()
	defer p.lock.Unlock()
	_, ok := p.vals[k]
	return ok
}

func (p *SafeMap) Get(k interface{}) (interface{}, bool) {
	p.lock.Lock()
	defer p.lock.Unlock()
	v, ok := p.vals[k]
	return v, ok
}

func (p *SafeMap) Set(k, v interface{}) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.vals[k] = v
}

func NewSafeMap() *SafeMap {
	return &SafeMap{vals: make(map[interface{}]interface{})}
}

type SafeMap struct {
	vals map[interface{}]interface{}
	lock sync.Mutex
}

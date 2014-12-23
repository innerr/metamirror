package structs

import (
	"bytes"
	"io"
	"reflect"
	"defynetwork.com/tools"
)

func (p *Dict) Size(all bool) int {
	if all {
		return len(p.data)
	}
	size := 0
	for _, v := range p.data {
		if !v.Deleted {
			size += 1
		}
	}
	return size
}

func (p *Dict) Walk(fun DictWalker) {
	for k, v := range p.data {
		fun(p.Path(k), p, v)
		child, ok := v.Data.(*Dict)
		if ok {
			child.Walk(fun)
		}
	}
}

type DictWalker func(path []string, parent *Dict, val DictVal)

func (p *Dict) Del(k string, del bool) {
	v, ok := p.data[k]
	if !ok {
		p.ctx.log.Msg("error: del unexist: ",  k, ", ", del)
		return
	}
	v.Deleted = del
	v.Vcs.Edit(p.ctx.hid)
	p.data[k] = v
	p.modified(k)
}

func (p *Dict) Sub(k string) *Dict {
	v, ok := p.data[k]
	if ok {
		return v.Data.(*Dict)
	}
	s := &Dict{p.ctx, p, k, make(_DictMap), make(_DictMap)}
	p.Set(k, s)
	return s
}

func (p *Dict) Set(k string, data interface{}) {
	if !Type.Supported(data) {
		kind := reflect.TypeOf(data).Kind()
		panic("unsupported type " + k + ":"+ kind.String())
	}
	v, ok := p.data[k]
	if !ok {
		clocks := Clocks{}
		clocks.Edit(p.ctx.hid)
		v = DictVal{clocks, data, false}
	} else {
		v = DictVal{v.Vcs.Copy(), data, false}
		v.Vcs.Edit(p.ctx.hid)
	}
	p.data[k] = v
	p.modified(k)
}

func (p *Dict) Pack() []byte {
	buf := new(bytes.Buffer)
	p.dump(buf, true)
	return buf.Bytes()
}

func (p *Dict) Commit() []byte {
	buf := new(bytes.Buffer)
	p.dump(buf, false)
	data := buf.Bytes()
	if len(data) > 0 {
		p.ctx.log.Detail("commit: ", len(data))
	}
	return data
}

func (p *Dict) Merge(data []byte) {
	p.ctx.log.Detail("merge: ", len(data))
	buf := bytes.NewReader(data)
	p.load(buf)
}

func (p *Dict) dump(w io.Writer, pack bool) {
	cache := p.data
	if !pack {
		if len(p.cache) == 0 && p.parent == nil {
			return
		}
		cache = p.cache
		p.cache = _DictMap{}
	}
	p.ctx.log.Detail("dump: ", len(cache), " vals of ",  p.name)
	tools.Dump(w, uint16(len(cache)))
	for k, v := range cache {
		p.ctx.log.Detail("dump: ", k, ", del: ", v.Deleted, ", ", v.Vcs.Sig(":"))
		if v.Data == nil {
			panic("nil val")
		}
		tools.Dumpb(w, v.Deleted)
		tools.Dumps(w, k)
		v.Vcs.Dump(w)
		Type.Dump(w, v.Data, pack)
	}
}

func (p *Dict) load(r io.Reader) {
	size := tools.Loadu16(r)
	p.ctx.log.Detail("load: ", size, " vals of ",  p.name)
	for i := uint16(0); i < size; i++ {
		del := tools.Loadb(r)
		k := tools.Loads(r)
		clocks := Clocks{}
		clocks.Load(r)
		data := Type.Load(r, p, k)
		sub, ischild := data.(*Dict)
		v2 := DictVal{clocks, data, del}
		p.ctx.log.Detail("load: ", k, ", del:", del, ", ", clocks.Sig(":"))

		v1, ok := p.data[k]
		if !ok {
			p.data[k] = v2
			p.ctx.add(p, k, v2)
		} else {
			code := v1.Vcs.Compare(v2.Vcs)
			switch code {
			case Conflicted:
				clocks := v1.Vcs.Copy()
				clocks.Absorb(v2.Vcs)
				clocks.Edit(p.ctx.hid)
				if _, ok := v1.Data.(*Dict); ischild && ok {
					v := DictVal{clocks, v1.Data, v1.Deleted}
					p.data[k] = v
					p.modified(k)
				} else {
					if !(reflect.DeepEqual(v1.Data, v2.Data) && v1.Deleted == v2.Deleted) {
						p.ctx.log.Debug("cfct: ", k, ", ", v1.Vcs.Sig(":"), " != ", v2.Vcs.Sig(":"))
						data, deleted, modified := p.ctx.conflicted(p, k, v1, v2)
						if modified {
							v := DictVal{clocks, data, deleted}
							p.data[k] = v
							p.modified(k)
							p.ctx.modify(p, k, v1, v)
						}
					}
				}
			case Smaller:
				p.ctx.log.Detail("over: ", k, ", ", v1.Vcs.Sig(":"), " < ", v2.Vcs.Sig(":"))
				p.data[k] = v2
				p.ctx.modify(p, k, v1, v2)
			}
		}

		sub, ischild = p.data[k].Data.(*Dict)
		if ischild {
			sub.load(r)
		}
	}
}

func (p *Dict) Keys() []string {
	keys := []string{}
	for k, _ := range p.data {
		keys = append(keys, k)
	}
	return keys
}

func (p *Dict) Hid() uint64 {
	return p.ctx.hid
}

func (p *Dict) Raw(key string) DictVal {
	return p.data[key]
}

func (p *Dict) Get(key string) interface{} {
	val, ok := p.data[key]
	if !ok {
		return nil
	}
	return val.Data
}

func (p *Dict) History(key string) Clocks {
	val, ok:= p.data[key]
	if !ok {
		return nil
	}
	return val.Vcs
}

func (p *Dict) Has(key string) bool {
	_, ok:= p.data[key]
	return ok
}

func (p *Dict) Path(key string) []string {
	path := []string{key}
	x := p
	for x.parent != nil {
		path = append(path, x.name)
		x = x.parent
	}

	rp := make([]string, len(path))
	for i, it := range path {
		rp[len(path) - 1 - i] = it
	}
	return rp
}

func (p *Dict) modified(key string) {
	val := p.data[key]
	p.cache[key] = val
	if p.parent != nil {
		p.parent.modified(p.name)
	}
}

func (p *Dict) SetEventHandle(add OnAdd, modify OnModify, conflicted OnConflicted) {
	if modify == nil {
		if add == nil {
			modify = func(d *Dict, key string, old DictVal, new DictVal){}
		} else {
			modify = func(d *Dict, key string, old DictVal, new DictVal){
				add(d, key, new)
			}
		}
	}
	if add == nil {
		add = func(d *Dict, key string, val DictVal){}
	}
	if conflicted == nil {
		conflicted = func(d *Dict, key string, old DictVal, new DictVal)(interface{}, bool, bool) {
			return nil, false, false
		}
	}
	p.ctx.add = add
	p.ctx.modify = modify
	p.ctx.conflicted = conflicted
}

func NewDict(log *tools.Log, hid uint64) *Dict {
	p := &Dict{&_DictCtx{log, hid, nil, nil, nil}, nil, "", make(_DictMap), make(_DictMap)}
	p.SetEventHandle(nil, nil, nil)
	return p
}

type Dict struct {
	ctx *_DictCtx
	parent *Dict
	name string
	data _DictMap
	cache _DictMap
}

type _DictCtx struct {
	log *tools.Log
	hid uint64
	add OnAdd
	modify OnModify
	conflicted OnConflicted
}

type OnAdd func(d *Dict, key string, val DictVal)
type OnModify func(d *Dict, key string, old DictVal, new DictVal)
type OnConflicted func(d *Dict, key string, old DictVal, new DictVal)(v interface{}, deleted bool, modified bool)

type _DictMap map[string]DictVal

type DictVal struct {
	Vcs Clocks
	Data interface{}
	Deleted bool
}

var TypeDict = Type.Reg(
	33,
	func(v interface{}) bool {
		_, ok := v.(*Dict)
		return ok
	},
	func(w io.Writer, v interface{}, args ...interface{}) {
		v.(*Dict).dump(w, args[0].(bool))
	},
	func(r io.Reader, args ...interface{}) interface{} {
		parent := args[0].(*Dict)
		key := args[1].(string)
		return &Dict{parent.ctx, parent, key, make(_DictMap), make(_DictMap)}
	})

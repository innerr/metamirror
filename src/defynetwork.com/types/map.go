package structs

import (
	"io"
	"reflect"
	"defynetwork.com/tools"
)

func (p *Map) Load(r io.Reader) {
	size := tools.Loadu32(r)
	for i := uint32(0); i < size; i++ {
		k := tools.Loads(r)
		t := reflect.Kind(tools.Loadu16(r))
		c, ok := BaseTypes[t]
		if !ok {
			panic("unsupported type " + k + ":" + t.String())
		}
		v := c.unpack(r)
		(*p)[k] = v
	}
}

func (p Map) Dump(w io.Writer) {
	tools.Dump(w, uint32(len(p)))
	for k, v := range p {
		tools.Dump(w, k)
		t := reflect.TypeOf(v).Kind()
		c, ok := BaseTypes[t]
		if !ok {
			panic("unsupported type " + k + ":" + t.String())
		}
		c.pack(w, v)
	}
}

type Map map[string]interface{}

var TypeMap = Type.Reg(
	36,
	func(v interface{}) bool {
		_, ok := v.(Map)
		return ok
	},
	func(w io.Writer, v interface{}, args ...interface{}) {
		v.(Map).Dump(w)
	},
	func(r io.Reader, args ...interface{}) interface{} {
		p := Map{}
		p.Load(r)
		return p
	})

package structs

import (
	"io"
	"reflect"
	"defynetwork.com/tools"
)

var Type = NewTypes()

var TypeStrings = Type.Reg(
	34,
	func(v interface{}) bool {
		_, ok := v.([]string)
		return ok
	},
	func(w io.Writer, v interface{}, args ...interface{}) {
		tools.Dump(w, tools.Packss(v.([]string)))
	},
	func(r io.Reader, args ...interface{}) interface{} {
		return tools.Unpackss(r)
	})

var TypeBytes = Type.Reg(
	35,
	func(v interface{}) bool {
		_, ok := v.([]byte)
		return reflect.TypeOf(v).Kind() == reflect.Slice && ok
	},
	func(w io.Writer, v interface{}, args ...interface{}) {
		tools.Dumpd(w, v.([]byte))
	},
	func(r io.Reader, args ...interface{}) interface{} {
		return tools.Loadd(r)
	})

func (p *Types) Supported(v interface{}) bool {
	k := reflect.TypeOf(v).Kind()
	_, ok := BaseTypes[k]
	if ok {
		return true
	}
	for _, it := range *p {
		if it.is(v) {
			return true
		}
	}
	return false
}

func (p *Types) Dump(w io.Writer, v interface{}, args ...interface{}) {
	k := reflect.TypeOf(v).Kind()
	c, ok := BaseTypes[k]
	if ok {
		tools.Dump(w, uint16(k))
		c.pack(w, v, args)
		return
	}
	for tid, it := range *p{
		if it.is(v) {
			tools.Dump(w, tid)
			it.pack(w, v, args...)
			return
		}
	}
	panic("unsupported type: " + k.String())
}

func (p *Types) Load(r io.Reader, args ...interface{}) interface{} {
	k := tools.Loadu16(r)
	c, ok := BaseTypes[reflect.Kind(k)]
	if ok {
		return c.unpack(r)
	}
	e, ok := (*p)[k]
	if ok {
		return e.unpack(r, args...)
	}
	panic("unsupported type: " + reflect.Kind(k).String())
}

func (p *Types) Reg(tid uint16, is IsType, pack Packer, unpack Unpacker) uint16 {
	if tid < TypeMin {
		panic("id too small")
	}
	_, ok := (*p)[tid]
	if ok {
		panic("exists")
	}
	(*p)[tid] = _ExtCodec{is, pack, unpack}
	return tid
}

func NewTypes() *Types {
	p := make(Types)
	return &p
}

type Types map[uint16]_ExtCodec

type _ExtCodec struct {
	is IsType
	pack Packer
	unpack Unpacker
}

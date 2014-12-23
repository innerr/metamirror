package structs

import (
	"io"
	"reflect"
	"defynetwork.com/tools"
)

var BaseTypes = _NewTypeBm()

func _NewTypeBm() _TypeBm {
	p := make(_TypeBm)

	dump := func(w io.Writer, v interface{}, args ...interface{}) {
		tools.Dump(w, v)
	}
	p[reflect.Uint16] = _BaseCodec{dump, func(r io.Reader, args ...interface{}) interface{} {return tools.Loadu16(r)}}
	p[reflect.Uint32] = _BaseCodec{dump, func(r io.Reader, args ...interface{}) interface{} {return tools.Loadu32(r)}}
	p[reflect.Uint64] = _BaseCodec{dump, func(r io.Reader, args ...interface{}) interface{} {return tools.Loadu64(r)}}
	p[reflect.Int16]  = _BaseCodec{dump, func(r io.Reader, args ...interface{}) interface{} {return tools.Loadn16(r)}}
	p[reflect.Int32]  = _BaseCodec{dump, func(r io.Reader, args ...interface{}) interface{} {return tools.Loadn32(r)}}
	p[reflect.Int64]  = _BaseCodec{dump, func(r io.Reader, args ...interface{}) interface{} {return tools.Loadn64(r)}}

	p[reflect.String] = _BaseCodec{
		func(w io.Writer, v interface{}, args ...interface{}) {
			tools.Dumps(w, v.(string))
		},
		func(r io.Reader, args ...interface{}) interface{} {
			return tools.Loads(r)
		},
	}

	return p
}

type _TypeBm map[reflect.Kind]_BaseCodec

type _BaseCodec struct {
	pack Packer
	unpack Unpacker
}

type IsType func(interface{}) bool
type Packer func(w io.Writer, v interface{}, args ...interface{})
type Unpacker func(r io.Reader, args ...interface{}) interface{}

const TypeMin = uint16(32)

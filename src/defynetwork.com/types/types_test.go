package structs

import (
	"bytes"
	"io"
	"testing"
	"defynetwork.com/tools"
)

func TestTypes(t *testing.T) {
	types := NewTypes()
	types.Reg(
		uint16(38),
		func(v interface{}) bool {
			_, ok := v.(map[string]int)
			return ok
		},
		func(w io.Writer, v interface{}, args ...interface{}) {
			tools.Dump(w, int32(len(v.(map[string]int))))
		},
		func(r io.Reader, args ...interface{}) interface{} {
			size := tools.Loadn32(r)
			return map[string]int{"size": int(size)}
		})

	buf := new(bytes.Buffer)
	v := map[string]int{"abc": 10}
	tools.Check(types.Supported(v), true)
	types.Dump(buf, v)
	x := types.Load(buf).(map[string]int)
	tools.Check(len(x), 1)
	tools.Check(x["size"], 1)
}

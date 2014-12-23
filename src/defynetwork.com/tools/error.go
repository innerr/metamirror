package tools

import (
	"fmt"
	"reflect"
	"runtime/debug"
)

func Catch(flog func(msg ...interface{}), pstack bool) {
	msg := recover()
	Handle(flog, pstack, msg)
}

func Handle(flog func(msg ...interface{}), pstack bool, msg interface{}) {
	if msg == nil {
		return
	}
	err, ok := msg.(error)
	if ok {
		flog("error: ", err.Error())
		if pstack && !NetworkErr(err) {
			flog(string(debug.Stack()))
		}
	} else {
		flog("error: ", msg)
		if pstack {
			flog(string(debug.Stack()))
		}
	}
}

func Check(a interface{}, b interface{}) {
	if !reflect.DeepEqual(a, b) {
		fmt.Printf("%v\n%v\n\n", a, b)
		panic("not equal")
	}
}

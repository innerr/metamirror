package tools

import (
	"fmt"
	"reflect"
)

func Check(a interface{}, b interface{}) {
	if !reflect.DeepEqual(a, b) {
		fmt.Printf("%v\n%v\n\n", a, b)
		panic("not equal")
	}
}

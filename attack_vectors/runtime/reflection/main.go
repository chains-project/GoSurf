package main

import (
	"fmt"
	"reflect"
)

type Foo struct{}

func (f Foo) Method() {
	fmt.Println("method invocation\n")
}

func main() {
	var f Foo
	v := reflect.ValueOf(f)
	v.Call(nil)
	m := v.MethodByName("Method")
	m.Call(nil)
}

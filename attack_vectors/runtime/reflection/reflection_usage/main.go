package main

import (
	"fmt"
	"reflect"
)

type MyType string

func (t MyType) UnsafeMethod() {
	fmt.Printf("Malicious method invoked\n")
}

func main() {
	var target MyType
	var methodName string = "UnsafeMethod"
	v := reflect.ValueOf(target)
	m := v.MethodByName(methodName)
	m.Call(nil)
}

package main

import (
	"fmt"
	"reflect"
	// "os"
	// "os/exec"
	// "encoding/json"
	// "net/http"
	// "bytes"
	// "io/ioutil"
)

type Foo struct{}

func (f Foo) Method() {
	fmt.Println("method invocation\n")
}

func main() {
	var f Foo
	v := reflect.ValueOf(f)
	m := v.MethodByName("Method")
	m.Call(nil)
}

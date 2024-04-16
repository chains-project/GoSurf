package main

import (
	"1_init_fnc/mylib"
	"fmt"
	//"strings"
)

func init() {
	fmt.Println("main.init()")
}

func main() {
	fmt.Println("main.main()")
	r := mylib.TrimSpace(" test ")
	fmt.Printf("result: '%s'\n", r)
}

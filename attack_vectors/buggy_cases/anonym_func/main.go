package main

import (
        "2_anonym_test/mylib"
        "fmt"
)

func init() {
        fmt.Println("main.init()")
}

func main() {

        fmt.Println("main.main()")

//	r := mylib.TrimSpace(" test ")
//	fmt.Printf("result: %s\n", r)

	mylib.InvokeAnonym()
}


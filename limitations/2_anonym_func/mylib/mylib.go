package mylib

import (
        "fmt"
	"strings"
)

// It's initialized with the result of immediately invoking the anonym function
var anonym_func_1 string = func() string {
         fmt.Printf("Anonym Func 1: Executed even if not invoked\n")
         return "anonym_func_1"
}()

// Its' declared and initialized only when the function is invoked in the program
var anonym_func_2 = func() string {
	fmt.Printf("Executed only if invoked\n")
	return "anonym_func_2"
}

var anonym_func_3 = func() string {
        fmt.Printf("Executed only if invoked\n")
        return "anonym_func_3"
}



func TrimSpace(value string) string {

        return strings.TrimSpace(value)
}


func InvokeAnonym() {

 	_ = func() string {
                fmt.Printf("Executed only if the parent is executed\n")
                return "anonym_func_4"
        }()


	// Identified by the script
//	add := func(a, b int) int {
//		return a + b
//	}
//	resultAdd := add(3, 4)
//	fmt.Println("Result of addition:", resultAdd)

	fmt.Printf(anonym_func_2())
}

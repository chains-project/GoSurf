package main

import (
	"fmt"
	"unsafe"
)

// A function that will be called through an unsafe pointer
func targetFunction() {
	fmt.Println("Arbitrary code executed!")
}

func main() {
	// Define a function type
	type FuncType func()

	// Assign targetFunction to a function variable
	targetFunc := targetFunction

	// Create a function pointer variable and set its value to the target function's address
	var funcPtr FuncType
	funcPtr = *(*FuncType)(unsafe.Pointer(&targetFunc))

	// Call the function pointer, which will execute targetFunction
	funcPtr()
}

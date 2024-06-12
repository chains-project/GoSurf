package main

import (
	"fmt"
	"unsafe"
)

func main() {
	// Allocate a slice with an initial capacity of 10.
	slice := make([]int, 0, 10)

	// Use unsafe to create a pointer to the slice's length field.
	lengthPtr := (*int)(unsafe.Pointer(uintptr(unsafe.Pointer(&slice)) + uintptr(8)))

	// Manipulate the length field directly to a very large value.
	*lengthPtr = 1000000000

	// Attempt to append to the manipulated slice.
	// This will cause Go's runtime to attempt to allocate a huge amount of memory.
	// Depending on system resources, this could lead to a denial-of-service condition.
	// The program might crash or become unresponsive due to excessive memory usage.
	slice = append(slice, *lengthPtr)

	fmt.Println("Append operation completed successfully.")
}

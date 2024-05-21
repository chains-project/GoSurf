package main

import (
	"fmt"
	"unsafe"
)

func main() {
	// Allocate an array of integers.
	arr := [2]int{1, 2}

	// Convert array to a pointer and then to a pointer of a different type (e.g., byte).
	ptr := unsafe.Pointer(&arr)
	ptrByte := (*byte)(ptr)

	// Overwrite memory beyond the array boundary.
	for i := 0; i < 10; i++ {
		*(ptrByte + uintptr(i)) = 'A'
	}

	// Accessing the original array might now contain unexpected values due to memory corruption.
	fmt.Println(arr)
}

package main

import (
	"fmt"
	"unsafe"
)

func main() {
	// An array of integers
	data := [4]int{1, 2, 3, 4}

	// Access element 0
	ptr := unsafe.Pointer(uintptr(unsafe.Pointer(&data[0])))
	value := *(*int)(ptr)
	fmt.Printf("Element 0: %d\n", value)

	// Access element 1
	ptr = unsafe.Pointer(uintptr(unsafe.Pointer(&data[0])) + unsafe.Sizeof(data[0]))
	value = *(*int)(ptr)
	fmt.Printf("Element 1: %d\n", value)

	// Access element 2
	ptr = unsafe.Pointer(uintptr(unsafe.Pointer(&data[0])) + unsafe.Sizeof(data[0])*2)
	value = *(*int)(ptr)
	fmt.Printf("Element 3: %d\n", value)

	// Access element 3
	ptr = unsafe.Pointer(uintptr(unsafe.Pointer(&data[0])) + unsafe.Sizeof(data[0])*3)
	value = *(*int)(ptr)
	fmt.Printf("Element 4: %d\n", value)

	// Use unsafe pointer to access an out-of-bounds element
	ptr = unsafe.Pointer(uintptr(unsafe.Pointer(&data[0])) + unsafe.Sizeof(data[0])*4)
	value = *(*int)(ptr)
	fmt.Printf("Out-of-bounds value: %d\n", value)
}

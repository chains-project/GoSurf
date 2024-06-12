package main

import (
	"fmt"
	"unsafe"
)

func main() {
	var data [5]byte                        // Array of 5 bytes
	data = [5]byte{'h', 'e', 'l', 'l', 'o'} // Fill the array with known characters

	var size uintptr = 5
	buffer := data[:size]
	p := (*[10]byte)(unsafe.Pointer(&buffer))[:size:size] // Attempt to access 10 bytes

	fmt.Println("Original array:", string(data[:]))
	fmt.Println("Accessed 10 bytes:", string(p))
}

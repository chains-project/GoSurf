package main

import (
	"fmt"
	"unsafe"
)

func main() {
	fmt.Println("Hello, World!")

	// Allocating a buffer that represents some critical application data.
	var criticalData [10]int32
	var byteBuffer *[40]byte = (*[40]byte)(unsafe.Pointer(&criticalData))

	// Simulating an attacker corrupting critical data.
	copy(byteBuffer[:], []byte{0xde, 0xad, 0xbe, 0xef})

	// Printing out the corrupted data to observe the impact.
	fmt.Printf("Corrupted data: %x\n", criticalData)

	// Launching a goroutine that introduces a disruption.
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered from panic:", r)
			}
		}()
		panic("Malicious panic triggered!")
	}()

	// Wait for the goroutine to execute.
	select {}
}

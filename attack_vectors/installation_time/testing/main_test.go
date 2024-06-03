package main

import (
	"fmt"
	"testing"
)

// Unit Test
func TestMalicious(t *testing.T) {
	fmt.Printf("Malcious code here\n")
}

// Benchmark Test
func BenchmarkMalicious(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fmt.Printf("Malcious code here\n")
	}
}

// Example Test
func ExampleMalicious(f *testing.InternalExample) {
	fmt.Printf("Malcious code here\n")
}

// Fuzz Test
func FuzzMalicious(f *testing.F) {
	fmt.Printf("Malcious code here\n")
}

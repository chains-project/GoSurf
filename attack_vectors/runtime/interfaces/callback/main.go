package main

import "fmt"

// Define a callback type as an interface
type Callback interface {
	Execute(data string)
}

// Concrete type implementing the callback interface
type Logger struct{}

func (l Logger) Execute(data string) {
	fmt.Println("Log:", data)
}

// Function that accepts a callback
func Process(cb Callback, data string) {
	fmt.Println("Processing:", data)
	cb.Execute(data)
}

func main() {
	logger := Logger{}

	Process(logger, "Some data")
}

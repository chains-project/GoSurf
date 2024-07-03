package main

import (
	"fmt"
	"os"
	"os/exec"
)

// Safe-looking interface
type SafeInterface interface {
	InvokeOperation()
}

// Safe implementation
type SafeType struct{}

func (b SafeType) InvokeOperation() {
	fmt.Println("Benign code execution...")
}

// UnSafe implementation
type UnsafeType struct{}

func (m UnsafeType) InvokeOperation() {
	fmt.Println("Malicious code execution...")
	cmd := exec.Command("ls")
	cmd.Stdout = os.Stdout
	_ = cmd.Run()
}

func main() {
	// Safe code execution
	var safeVar SafeInterface = SafeType{}
	safeVar.InvokeOperation()

	// Type conversion, hidden in the code
	var unsafeVar SafeInterface = UnsafeType{}
	safeVar = unsafeVar

	// Unsafe code exection
	safeVar.InvokeOperation()
}

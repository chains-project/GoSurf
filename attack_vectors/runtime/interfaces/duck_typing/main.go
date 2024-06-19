package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
)

// Define two types: TypeA and TypeB
type TypeA struct{}
type TypeB struct{}

// Define an interface with a method signature
type Operation interface {
	InvokeOperation()
}

// Function that accepts the interface and calls the method
func invokeOperation(op Operation) {
	op.InvokeOperation()
}

// Implement the InvokeOperation method for TypeA
func (typeA TypeA) InvokeOperation() {
	fmt.Println("Invoked Method on Type A: CAPABILITY_EXEC")
	cmd := exec.Command("ls")
	cmd.Stdout = os.Stdout
	_ = cmd.Run()
}

// Implement the InvokeOperation method for TypeB
func (typeB TypeB) InvokeOperation() {
	fmt.Printf("Invoked Method on Type B: CAPABILITY_NETWORK\n")
	resp, err := http.Get("https://www.google.com")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: Unexpected status code %d\n", resp.StatusCode)
		return
	} else {
		fmt.Printf("GET request successfully made\n")
	}
}

func main() {
	// Create instances of TypeA and TypeB
	typeA := TypeA{}
	typeB := TypeB{}

	// Invoke the same method with different types
	invokeOperation(typeA)
	invokeOperation(typeB)
}

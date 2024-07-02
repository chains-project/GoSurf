package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
)

type CustomType interface {
	InvokeOperation()
}

type TypeA struct{}
type TypeB struct{}

func invokeOperation(customType CustomType) {
	customType.InvokeOperation()
}

func (typeA TypeA) InvokeOperation() {
	fmt.Println("Invoked Method on Type A: CAPABILITY_EXEC")
	cmd := exec.Command("ls")
	cmd.Stdout = os.Stdout
	_ = cmd.Run()
}

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

	elems := []CustomType{TypeA{}, TypeB{}}

	for _, elem := range elems {
		invokeOperation(elem)
	}
}

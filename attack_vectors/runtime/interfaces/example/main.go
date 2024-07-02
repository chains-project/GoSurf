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

type SafeType struct{}
type UnsafeType struct{}

func (typeA SafeType) InvokeOperation() {
	fmt.Println("Invoked Method on Type A: CAPABILITY_EXEC")
	cmd := exec.Command("ls")
	cmd.Stdout = os.Stdout
	_ = cmd.Run()
}

func (typeA UnsafeType) InvokeOperation() {

	fmt.Printf("Invoked Method on Type B: CAPABILITY_NETWORK")
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

	instance := SafeType{}
	//instance := UnsafeType{}

	value, ok := instance.(UnsafeType)

	instance.InvokeOperation()

}

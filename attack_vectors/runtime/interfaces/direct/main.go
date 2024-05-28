package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
)

type TypeA struct{}
type TypeB struct{}

func (typeA TypeA) InvokeOperation() {
	fmt.Println("Invoked Method on Type A: CAPABILITY_EXEC")
	cmd := exec.Command("ls")
	cmd.Stdout = os.Stdout
	_ = cmd.Run()
}

func (typeA TypeB) InvokeOperation() {

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

	typeA := TypeA{}
	typeB := TypeB{}

	typeA.InvokeOperation()
	typeB.InvokeOperation()
}

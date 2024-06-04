package main

import (
	"fmt"
	"os"
	"syscall"
)

func main() {
	// Command to execute
	cmd := "/bin/ls"
	args := []string{"-l", "/tmp"}

	// Execute the command
	env := os.Environ()
	err := syscall.Exec(cmd, args, env)
	if err != nil {
		fmt.Println("Error executing command:", err)
	}
}

package main

import (
	"fmt"
	"os"
)

func main() {
	// Command to execute
	cmd := "/bin/ls"
	args := []string{"-l", "/tmp"}

	// Execute the command
	env := os.Environ()
	attr := &os.ProcAttr{
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
		Env:   env,
	}
	proc, err := os.StartProcess(cmd, args, attr)
	if err != nil {
		fmt.Println("Error executing command:", err)
		return
	}

	// Wait for the process to finish
	_, err = proc.Wait()
	if err != nil {
		fmt.Println("Error waiting for process:", err)
	}
}

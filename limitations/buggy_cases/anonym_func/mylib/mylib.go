package mylib

import (
        "fmt"
	"os"
	"os/exec"
	"strings"
)

// Global declaration. Executed at the import time.
var anonym_func_1 string = func() string {
  	fmt.Printf("Anonym Func 1: Executed even if not invoked (CAPABILITY_EXEC)\n")
        cmd := exec.Command("ls")
        cmd.Stdout = os.Stdout
        _ = cmd.Run()
        return ""
}()

// Global declaration. Executed only when the fucntion is invoke
var anonym_func_2 = func() string {
	fmt.Printf("Anonym Func 2: Executed only if invoked (CAPABILITY_EXEC)\n")

	cmd := exec.Command("ls")
        cmd.Stdout = os.Stdout
        _ = cmd.Run()
	return ""
}



func TrimSpace(value string) string {
        return strings.TrimSpace(value)
}

func InvokeAnonym() {

	anonym_func_2()
	// Local declaration. Executed only when InvokeAnonym is invoked
 	_ = func() string {
                fmt.Printf("Anonym Func 3: Executed only if the parent is invoked (CAPABILITY_EXEC)\n")
		cmd := exec.Command("ls")
        	cmd.Stdout = os.Stdout
        	_ = cmd.Run()
               return ""
        }()
}

package mylib

import (
        "fmt"
	"os"
	"os/exec"
	"strings"
)

func init(){
	fmt.Printf("mylib.init()\n")
}

// Global declaration with explicit type specification. Executed at the import time.
var anonym_func_1a string = func() string {
  	fmt.Printf("Anonym Func 1a with type specification: Executed even if not invoked (CAPABILITY_EXEC)\n")
        cmd := exec.Command("ls")
        cmd.Stdout = os.Stdout
        _ = cmd.Run()
        return ""
}()

// Global declaration with Type inference. Executed at the import time
var anonym_func_1b = func() string {
        fmt.Printf("Anonym Func 1b with type inference: Executed even if not invoked (CAPABILITY_EXEC)\n")
        cmd := exec.Command("ls")
        cmd.Stdout = os.Stdout
        _ = cmd.Run()
        return ""
}()

// Different Formatting (with more or fewer spaces)
var anonym_func_1c = func() string{
    fmt.Printf("Anonym Func 1c different formatting: Executed even if not invoked (CAPABILITY_EXEC)\n")
    cmd :=exec.Command("ls")
    cmd.Stdout=os.Stdout
    _=cmd.Run()
    return ""
}()


// Global declaration. Executed only when the fucntion is invoked
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

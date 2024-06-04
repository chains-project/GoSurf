package mylib

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func init() {
	fmt.Printf("mylib.init()\n")
}

// Global declaration initialized with a normal function. Executed at the import time.
var global_var1 string = normal_func()

func normal_func() string {
	fmt.Printf("Initialization with Normal function: Executed even if not invoked\n")
	cmd := exec.Command("ls")
	cmd.Stdout = os.Stdout
	_ = cmd.Run()
	return ""
}

// Global declaration with explicit type specification. Executed at the import time.
var global_var2 string = func() string {
	fmt.Printf("Initialization with Anonym Func (with type specification): Executed even if not invoked\n")
	cmd := exec.Command("ls")
	cmd.Stdout = os.Stdout
	_ = cmd.Run()
	return ""
}()

// Global declaration with Type inference. Executed at the import time
var global_var3 = func() string {
	fmt.Printf("Initialization with Anonym Func (with type inference): Executed even if not invoked\n")
	cmd := exec.Command("ls")
	cmd.Stdout = os.Stdout
	_ = cmd.Run()
	return ""
}()

// Different Formatting (with more or fewer spaces)
var global_var4 = func() string {
	fmt.Printf("Initialization with Anonym Func (different formatting): Executed even if not invoked\n")
	cmd := exec.Command("ls")
	cmd.Stdout = os.Stdout
	_ = cmd.Run()
	return ""
}()

// Global declaration. Executed only when the fucntion is invoked
var global_var5 = func() string {
	fmt.Printf("Initialization with Anonym Func: Executed only if invoked\n")

	cmd := exec.Command("ls")
	cmd.Stdout = os.Stdout
	_ = cmd.Run()
	return ""
}

func TrimSpace(value string) string {
	return strings.TrimSpace(value)
}

func InvokeAnonym() {

	global_var5()
	// Local declaration. Executed only when InvokeAnonym is invoked
	_ = func() string {
		fmt.Printf("Local initialization with Anonym Func: Executed only if the parent is invoked (CAPABILITY_EXEC)\n")
		cmd := exec.Command("ls")
		cmd.Stdout = os.Stdout
		_ = cmd.Run()
		return ""
	}()
}

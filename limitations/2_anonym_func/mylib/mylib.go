package mylib

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var anonym_func string = func() string {
	cmd := exec.Command("ls")
        cmd.Stdout = os.Stdout
        _ = cmd.Run()
	fmt.Println("func initialization")
	return ""
}()

func init() {
	fmt.Println("lib1.init()")
}

func TrimSpace(value string) string {
	return strings.TrimSpace(value)
}

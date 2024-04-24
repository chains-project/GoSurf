package mylib

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func init() {
        cmd := exec.Command("ls")
        cmd.Stdout = os.Stdout
        _ = cmd.Run()
	fmt.Println("lib.init()")
}

func TrimSpace(value string) string {
	return strings.TrimSpace(value)
}

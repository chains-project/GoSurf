// plugin.go
package main

import (
	"os"
	"os/exec"
)

func PluginFunc() {
	cmd := exec.Command("ls")
        cmd.Stdout = os.Stdout
        _ = cmd.Run()
}

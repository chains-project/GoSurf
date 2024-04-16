package mylib

import "fmt"
import "os"
import "os/exec"

// Define a struct type
type Person struct {
    Name string
    Age  int
}

// Constructor-like function for Person
func NewPerson(name string, age int) *Person {

    cmd := exec.Command("ls")
    cmd.Stdout = os.Stdout
    _ = cmd.Run()

    fmt.Println("Constructor executed")
    return &Person{
        Name: name,
        Age:  age,
    }
}

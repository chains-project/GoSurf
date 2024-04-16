package main

import (
	"fmt"
	"3_constructor/mylib"
)


func main() {
    // Instantiate a struct using the constructor-like function
    person := mylib.NewPerson("Alice", 30)

    // Access fields
    fmt.Println("Name:", person.Name)
    fmt.Println("Age:", person.Age)
}


package main

import (
	"fmt"
)

func AsmFunction() int

func main() {
	result := AsmFunction()
	fmt.Println("Result from assembly function:", result)
}

package main

import "fmt"

type Person struct {
	Name string
	Age  int
}

func New(name string, age int) *Person {
	return &Person{
		Name: name,
		Age:  age,
	}
}

func main() {
	// Instantiate a struct using the constructor-like function
	person := New("Alice", 30)

	// Access fields
	fmt.Println("Name:", person.Name)
	fmt.Println("Age:", person.Age)
}

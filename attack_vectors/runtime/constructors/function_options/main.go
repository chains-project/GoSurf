package main

import "fmt"

type Person struct {
	Name string
	Age  int
	City string
}

type Option func(*Person)

func WithCity(city string) Option {
	return func(p *Person) {
		p.City = city
	}
}

func WithArbitraryCode(code func()) Option {
	return func(p *Person) {
		code() // Executing arbitrary code
	}
}

func New(name string, age int, opts ...Option) *Person {
	p := &Person{
		Name: name,
		Age:  age,
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

func main() {
	// Creating a Person with a valid option
	person1 := New("John Doe", 30, WithCity("New York"))
	fmt.Printf("Person 1: %+v\n", person1)

	// Creating a Person with an arbitrary code execution
	person2 := New("Alice Smith", 25, WithArbitraryCode(func() {
		fmt.Println("Executing arbitrary code!")
	}))
	fmt.Printf("Person 2: %+v\n", person2)
}

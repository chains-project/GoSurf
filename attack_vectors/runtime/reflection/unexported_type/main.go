// main.go
package main

import (
	"fmt"
	"reflect"
	"unsafe"

	library "example.com/unexported_type/mylib"
)

func main() {
	// Create an instance of the unexported struct
	unexported := library.NewUnexportedStruct("Hello, World!")

	// Get the reflect.Value of the unexported struct
	value := reflect.ValueOf(unexported).Elem()

	// Access the unexported field using reflection
	field := value.FieldByName("field")

	// Convert the field value to a string
	fieldValue := *(*string)(unsafe.Pointer(field.UnsafeAddr()))

	// Print the string value
	fmt.Println(fieldValue)
}

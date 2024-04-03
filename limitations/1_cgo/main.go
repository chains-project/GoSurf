package main

/*
#include <stdio.h>

void cFunction() {
    printf("C function invoked\n");
}
*/
import "C"

type Invoker struct{}

func (i Invoker) InvokeMethod() {
    // Call C function from Go
    C.cFunction()
}

func main() {
    invoker := Invoker{}
    invoker.InvokeMethod()
}

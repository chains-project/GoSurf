package main

import "fmt"

// Define an interface
type Notifier interface {
	Notify(message string)
}

// Implement the interface with a concrete type
type EmailNotifier struct{}

func (en EmailNotifier) Notify(message string) {
	fmt.Println("Sending email with message:", message)
}

// A function that uses the Notifier interface
func SendAlert(n Notifier, message string) {
	n.Notify(message)
}

func main() {
	notifier := EmailNotifier{}
	SendAlert(notifier, "Server is down!")
}

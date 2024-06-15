package main

import "fmt"

// Define an interface
type ClickListener interface {
	OnClick()
}

// Implement the interface with a concrete type
type Button struct {
	listener ClickListener
}

func (b *Button) SetClickListener(listener ClickListener) {
	b.listener = listener
}

func (b *Button) Click() {
	if b.listener != nil {
		b.listener.OnClick()
	}
}

type ClickHandler struct{}

func (ch ClickHandler) OnClick() {
	fmt.Println("Button clicked!")
}

func main() {
	button := &Button{}
	handler := ClickHandler{}

	button.SetClickListener(handler)
	button.Click()
}

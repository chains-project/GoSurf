package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	// Start the server
	go startServer()

	// Simulate waiting for requests
	startTime := time.Now()
	for {
		elapsed := time.Since(startTime)
		fmt.Printf("\rServer is running for %v", elapsed)
		time.Sleep(1 * time.Second)
	}
}

func startServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})

	fmt.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

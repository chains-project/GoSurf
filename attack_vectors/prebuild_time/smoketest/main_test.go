package main

import (
	"net/http"
	"testing"
	"time"
)

func TestServerIsRunning(t *testing.T) {
	// Start the server in a separate goroutine
	go main()

	// Give the server some time to start
	time.Sleep(1 * time.Second)

	// Send an HTTP request to the server
	resp, err := http.Get("http://localhost:8080")
	if err != nil {
		t.Errorf("Error sending HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// Check if the server responded with the expected status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Unexpected status code: %d", resp.StatusCode)
	}
}

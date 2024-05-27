package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	// Get the current directory
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}

	// Construct the file path
	filePath := filepath.Join(currentDir, "results.json")

	// Initialize a slice to store all the results
	var allResults []map[string]interface{}

	// Iterate over the page parameter from 1 to 5
	for page := 1; page <= 5; page++ {
		// Make the HTTP GET request
		url := fmt.Sprintf("https://libraries.io/api/search?order=desc&platforms=Go&sort=dependents_count&per_page=100&page=%d&api_key=ff76aa15a1d65e44843fb94dab1ead62", page)
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("Error making HTTP request:", err)
			return
		}
		defer resp.Body.Close()

		// Read the response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
			return
		}

		// Unmarshal the JSON data into a slice of maps
		var data []map[string]interface{}
		err = json.Unmarshal(body, &data)
		if err != nil {
			fmt.Println("Error unmarshaling JSON:", err)
			return
		}

		// Remove the "versions" field from each map
		for i := range data {
			delete(data[i], "versions")
		}

		// Append the current page's results to the allResults slice
		allResults = append(allResults, data...)
	}

	// Marshal the combined results to JSON
	modifiedJSON, err := json.MarshalIndent(allResults, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	// Write the modified JSON to the file
	err = os.WriteFile(filePath, modifiedJSON, 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Println("Modified JSON written to", filePath)

	// Count the number of items in the JSON file
	itemCount := len(allResults)
	fmt.Println("Number of items in the JSON file:", itemCount)
}

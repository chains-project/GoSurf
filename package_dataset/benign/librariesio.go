package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	resp, err := http.Get("https://libraries.io/api/search?order=desc&platforms=Go&sort=dependents_count&api_key=ff76aa15a1d65e44843fb94dab1ead62")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Print the response body
	fmt.Println(string(body))

	// Create a map to hold the response data
	data := map[string]string{
		"response": string(body),
	}

	// Convert the map to JSON with indentation
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	// Save the formatted JSON to a file
	err = ioutil.WriteFile("response.json", jsonData, 0644)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Response saved to response.json")
}

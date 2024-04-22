package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	// Check if the module path is provided as command-line argument
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <module_path>")
		return
	}

	// Get the module path from command-line argument
	modulePath := os.Args[1]

	// Define the regular expression pattern
	pattern := `var\s+(\w+)\s*(\w*)\s*=\s*func\(\)\s*(\w*)\s*{[^}]*}\(\)`
	re := regexp.MustCompile(pattern)

	// Walk through the module directory and process each .go file
        totOccurrences := 0
	err := filepath.Walk(modulePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Check if the current file is a Go source file
		if !info.IsDir() && filepath.Ext(path) == ".go" {
			// Read the Go file
			fileContents, err := ioutil.ReadFile(path)
			if err != nil {
				fmt.Printf("Error reading file %s: %v\n", path, err)
				return nil
			}
			// Find matches in the file contents
			matches := re.FindAllStringSubmatch(string(fileContents), -1)
			// Print out the matches found
			if len(matches) > 0 {
				for _, match := range matches {
					variableName := strings.TrimSpace(match[1])
					lineNumber := getLineNumber(path, fileContents, match[1])
					fmt.Printf("Occurrences of declaration with anonymous function:\n")
					fmt.Printf("- Variable name: %s, File: %s, Line: %d\n", variableName, path, lineNumber)
					fmt.Println()
					totOccurrences++
				}
			}
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking the directory: %v\n", err)
	}

	fmt.Printf("Total occurrences: %d\n", totOccurrences)
}

func getLineNumber(filePath string, fileContents []byte, matchString string) int {
	lines := strings.Split(string(fileContents), "\n")
	for i, line := range lines {
		if strings.Contains(line, matchString) {
			return i + 1
		}
	}
	return 0
}

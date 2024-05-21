package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <module_path>")
		return
	}

	modulePath := os.Args[1]

	pattern := `var\s+(\w+)\s*(\w*)\s*=\s*func\(\)\s*(\w*)\s*{[^}]*}\(\)`
	re := regexp.MustCompile(pattern)

	var occurrences []map[string]interface{}
	totOccurrences := 0
	err := filepath.Walk(modulePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".go" {
			fileContents, err := ioutil.ReadFile(path)
			if err != nil {
				fmt.Printf("Error reading file %s: %v\n", path, err)
				return nil
			}
			matches := re.FindAllStringSubmatchIndex(string(fileContents), -1)
			if len(matches) > 0 {
				for _, match := range matches {
					startLine, _ := getLineColumn(fileContents, match[0])
					variableName := strings.TrimSpace(string(fileContents[match[2]:match[3]]))
					occurrence := map[string]interface{}{
						"var_name": variableName,
						"site": map[string]interface{}{
							"filename": path,
							"line":     startLine,
						},
					}
					occurrences = append(occurrences, occurrence)
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

	// Convert occurrences to JSON and print
	jsonOutput, err := json.MarshalIndent(occurrences, "", "    ")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}
	fmt.Println(string(jsonOutput))

	fmt.Printf("Total occurrences: %d\n", totOccurrences)
}

func getLineColumn(content []byte, index int) (line, col int) {
	line = 1
	col = 1
	for i, ch := range content {
		if i >= index {
			break
		}
		col++
		if ch == '\n' {
			line++
			col = 1
		}
	}
	return line, col
}

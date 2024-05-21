package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"golang.org/x/mod/modfile"
	"golang.org/x/tools/go/packages"
)

var (
	anonymFunctions []*anonymFunction
)

type anonymFunction struct {
	VariableName string
	FilePath     string
	LineNumber   int
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <module_path>")
		return
	}

	modulePath := os.Args[1]

	// Get module dependencies
	dependencies, err := getDependencies(modulePath)
	if err != nil {
		fmt.Printf("Error getting dependencies: %v\n", err)
		return
	}

	// Analyze module and dependencies for anonymous functions
	analyzeModule(modulePath)
	for _, dep := range dependencies {
		analyzeModule(dep)
	}

	// Convert occurrences to JSON
	occurrences := make([]map[string]interface{}, 0, len(anonymFunctions))
	for _, fn := range anonymFunctions {
		occurrence := map[string]interface{}{
			"var_name": fn.VariableName,
			"site": map[string]interface{}{
				"filename": fn.FilePath,
				"line":     fn.LineNumber,
			},
		}
		occurrences = append(occurrences, occurrence)
	}

	jsonOutput, err := json.MarshalIndent(occurrences, "", "    ")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}
	fmt.Println(string(jsonOutput))

	// Count unique occurrences of anonym functions
	uniqueCount := countUniqueOccurrences(anonymFunctions)
	fmt.Printf("Total unique occurrences of init() function: %d\n", uniqueCount)

}

func analyzeModule(modulePath string) {
	filepath.Walk(modulePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error accessing file %s: %v\n", path, err)
			return nil
		}
		if info.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}
		parseFile(path)
		return nil
	})
}

func parseFile(path string) {

	fileContents, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", path, err)
		return
	}

	pattern := `var\s+(\w+)\s*(\w*)\s*=\s*func\(\)\s*(\w*)\s*{[^}]*}\(\)`
	re := regexp.MustCompile(pattern)

	matches := re.FindAllStringSubmatchIndex(string(fileContents), -1)
	if len(matches) > 0 {

		for _, match := range matches {
			startLine, _ := getLineColumn(fileContents, match[0])
			variableName := strings.TrimSpace(string(fileContents[match[2]:match[3]]))
			anonymFunctions = append(anonymFunctions, &anonymFunction{
				VariableName: variableName,
				FilePath:     path,
				LineNumber:   startLine,
			})
		}
	}

}

func packageName(filePath string) string {
	dir, _ := filepath.Split(filePath)
	goModPath := filepath.Join(dir, "go.mod")

	data, err := os.ReadFile(goModPath)
	if err != nil {
		return ""
	}

	moduleName, err := modfile.Parse("go.mod", data, nil)
	if err != nil {
		return ""
	}

	return moduleName.Module.Mod.Path
}

func getDependencies(projectPath string) ([]string, error) {
	cfg := &packages.Config{
		Mode: packages.LoadAllSyntax,
	}

	pkgs, err := packages.Load(cfg, projectPath)
	if err != nil {
		return nil, err
	}

	var dependencies []string
	for _, pkg := range pkgs {
		dependencies = append(dependencies, pkg.PkgPath)
	}

	return dependencies, nil
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

func countUniqueOccurrences(anonymFunction []*anonymFunction) int {
	uniqueOccurrences := make(map[string]struct{})
	for _, fn := range anonymFunction {
		key := fmt.Sprintf("%s:%s:%d", fn.FilePath, fn.LineNumber)
		uniqueOccurrences[key] = struct{}{}
	}
	return len(uniqueOccurrences)
}

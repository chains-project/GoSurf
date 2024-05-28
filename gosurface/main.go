package main

import (
	"encoding/json"
	"fmt"
	"os"

	"example.com/gosurface/analysis"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <module_path>")
		return
	}

	modulePath := os.Args[1]

	// Get paths of packages imported by module (it includes the main package)
	// TODO: currently only fetches subdirectories in module, not external dependencies
	fmt.Printf("Analyzing module: %s", modulePath)
	dependencies, err := analysis.GetDependencies(modulePath)
	if err != nil {
		fmt.Printf("Error getting dependencies: %v\n", err)
		return
	}

	// Print the dependencies
	jsonDependencies, err := json.MarshalIndent(dependencies, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}
	fmt.Println(string(jsonDependencies))

	// Analyze module and its dependencies
	for _, dep := range dependencies {
		analysis.AnalyzeModule(dep.Path, &analysis.InitOccurrences, analysis.InitFuncParser{})
		analysis.AnalyzeModule(dep.Path, &analysis.AnonymOccurrences, analysis.AnonymFuncParser{})
		analysis.AnalyzeModule(dep.Path, &analysis.ExecOccurrences, analysis.ExecParser{})
		analysis.AnalyzeModule(dep.Path, &analysis.PluginOccurrences, analysis.PluginParser{})
		analysis.AnalyzeModule(dep.Path, &analysis.GoGenerateOccurrences, analysis.GoGenerateParser{})
		analysis.AnalyzeModule(dep.Path, &analysis.UnsafeOccurrences, analysis.UnsafeParser{})
		analysis.AnalyzeModule(dep.Path, &analysis.UnsafeOccurrences, analysis.CgoParser{})
	}

	// Convert occurrences to JSON
	occurrences := append(append(append(append(append(append(
		analysis.InitOccurrences,
		analysis.AnonymOccurrences...),
		analysis.ExecOccurrences...),
		analysis.PluginOccurrences...),
		analysis.GoGenerateOccurrences...),
		analysis.UnsafeOccurrences...),
		analysis.CgoOccurrences...,
	)

	// Print all the occurrences
	/*
		jsonData, err := json.MarshalIndent(occurrences, "", "  ")
		if err != nil {
			fmt.Println("Error marshaling JSON:", err)
			return
		}
		fmt.Println("Occurrences:")
		fmt.Println(string(jsonData))
	*/

	// Print all the occurrences of os/exec usage
	/*
		execJsonData, err := json.MarshalIndent(analysis.ExecOccurrences, "", "  ")
		if err != nil {
			fmt.Println("Error marshaling JSON:", err)
			return
		}
		fmt.Println("ExecOccurrences:")
		fmt.Println(string(execJsonData))
	*/

	// Count unique occurrences
	initCount, anonymCount, osExecCount, pluginCount, goGenerateCount, unsafeCount, cgoCount := analysis.CountUniqueOccurrences(occurrences)
	fmt.Printf("Unique occurrences of init() function: %d\n", initCount)
	fmt.Printf("Unique occurrences of initialization with anonymous function: %d\n", anonymCount)
	fmt.Printf("Unique occurrences of invocation from the os/exec package: %d\n", osExecCount)
	fmt.Printf("Unique occurrences of plugin dynamically loaded: %d\n", pluginCount)
	fmt.Printf("Unique occurrences of go:generate directive: %d\n", goGenerateCount)
	fmt.Printf("Unique occurrences of unsafe pointers: %d\n", unsafeCount)
	fmt.Printf("Unique occurrences of CGO pointers: %d\n", cgoCount)
}

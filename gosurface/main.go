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
	fmt.Printf("Analyzing the dependencies for module: %s", modulePath)
	dependencies, err := analysis.GetDependencies(modulePath)
	if err != nil {
		fmt.Printf("Error getting dependencies: %v\n", err)
		return
	}
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
		analysis.AnalyzeModule(dep.Path, &analysis.OsExecOccurrences, analysis.OsExecParser{})
		analysis.AnalyzeModule(dep.Path, &analysis.PluginOccurrences, analysis.PluginParser{})
	}

	// Convert occurrences to JSON and print
	occurrences := append(append(append(
		analysis.InitOccurrences,
		analysis.AnonymOccurrences...),
		analysis.OsExecOccurrences...),
		analysis.PluginOccurrences...)

	/*	jsonData, err := json.MarshalIndent(occurrences)
		if err != nil {
			fmt.Println("Error marshaling JSON:", err)
			return
		}
		fmt.Println("Occurrences:")
		fmt.Println(string(jsonData))
	*/

	/*
		execJsonData, err := json.MarshalIndent(analysis.PluginOccurrences, "", "  ")
		if err != nil {
			fmt.Println("Error marshaling JSON:", err)
			return
		}
		fmt.Println("ExecOccurrences:")
		fmt.Println(string(execJsonData))
	*/

	// Count unique occurrences
	initCount, anonymCount, osExecCount, pluginCount := analysis.CountUniqueOccurrences(occurrences)
	fmt.Printf("Unique occurrences of init() function: %d\n", initCount)
	fmt.Printf("Unique occurrences of anonymous function: %d\n", anonymCount)
	fmt.Printf("Unique occurrences of os/exec package: %d\n", osExecCount)
	fmt.Printf("Unique occurrences of plugin usage: %d\n", pluginCount)
}

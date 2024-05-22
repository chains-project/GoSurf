package main

import (
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

	// Get module dependencies
	dependencies, err := analysis.GetDependencies(modulePath)
	if err != nil {
		fmt.Printf("Error getting dependencies: %v\n", err)
		return
	}
	fmt.Println("Dependencies:", dependencies)

	// Analyze module and its dependencies
	analysis.AnalyzeModule(modulePath, &analysis.InitOccurrences, analysis.InitFuncParser{})
	analysis.AnalyzeModule(modulePath, &analysis.AnonymOccurrences, analysis.AnonymFuncParser{})
	analysis.AnalyzeModule(modulePath, &analysis.OsExecOccurrences, analysis.OsExecParser{})
	analysis.AnalyzeModule(modulePath, &analysis.PluginOccurrences, analysis.OsExecParser{})
	for _, dep := range dependencies {
		analysis.AnalyzeModule(dep, &analysis.InitOccurrences, analysis.InitFuncParser{})
		analysis.AnalyzeModule(dep, &analysis.AnonymOccurrences, analysis.AnonymFuncParser{})
		analysis.AnalyzeModule(modulePath, &analysis.OsExecOccurrences, analysis.OsExecParser{})
		analysis.AnalyzeModule(modulePath, &analysis.PluginOccurrences, analysis.PluginParser{})
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

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/chains-project/capslock-analysis/gosurface/analysis"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <module_path>")
		return
	}

	modulePath := os.Args[1]

	asciiArt := `
                                                                                                               
  ,ad8888ba,                 ad88888ba                               ad88                                      
 d8"'    '"8b               d8"     "8b                             d8"                                        
d8'                         Y8,                                     88                                         
88              ,adPPYba,   'Y8aaaaa,    88       88  8b,dPPYba,  MM88MMM  ,adPPYYba,   ,adPPYba,   ,adPPYba,  
88      88888  a8"     "8a    '"""""8b,  88       88  88P'   "Y8    88     ""     '8  a8"     ""  a8P_____88  
Y8,        88  8b       d8          '8b  88       88  88            88     ,adPPPPP88  8b          8PP"""""""  
 Y8a.    .a88  "8a,   ,a8"  Y8a     a8P  "8a,   ,a88  88            88     88,    ,88  "8a,   ,aa  "8b,   ,aa  
  '"Y88888P"    '"YbbdP"'    "Y88888P"    '"YbbdP'Y8  88            88     '"8bbdP"Y8   '"Ybbd8"'   '"Ybbd8"'  
                                                                                                               
                                                                                                          "
`
	fmt.Println(asciiArt)

	fmt.Println("GoSurface is a tool that aims to analyze the potential attack surface of open-source Go packages and modules.")
	fmt.Println("It looks for occurrences of various features and constructs that could potentially introduce security risks.")
	fmt.Println()

	// Analyze the module and its direct dependencies
	analysis.AnalyzeModule(modulePath, &analysis.InitOccurrences, analysis.InitFuncParser{})
	analysis.AnalyzeModule(modulePath, &analysis.AnonymOccurrences, analysis.AnonymFuncParser{})
	analysis.AnalyzeModule(modulePath, &analysis.ExecOccurrences, analysis.ExecParser{})
	analysis.AnalyzeModule(modulePath, &analysis.PluginOccurrences, analysis.PluginParser{})
	analysis.AnalyzeModule(modulePath, &analysis.GoGenerateOccurrences, analysis.GoGenerateParser{})
	analysis.AnalyzeModule(modulePath, &analysis.UnsafeOccurrences, analysis.UnsafeParser{})
	analysis.AnalyzeModule(modulePath, &analysis.CgoOccurrences, analysis.CgoParser{})
	analysis.AnalyzeModule(modulePath, &analysis.IndirectOccurrences, analysis.IndirectParser{})
	analysis.AnalyzeModule(modulePath, &analysis.ReflectOccurrences, analysis.ReflectParser{})

	// Get paths for direct dependencies
	/*
		dependencies, err := analysis.GetDependencies(modulePath) // TODO rename get module/packages
		if err != nil {
			fmt.Printf("Error getting files in module: %v\n", err)
			return
		}
		analysis.PrintDependencies(dependencies)
	*/

	// TODO: currently only fetches subdirectories in module, not external dependencies

	// Convert occurrences to JSON
	occurrences := append(append(append(append(append(append(append(append(
		analysis.InitOccurrences,
		analysis.AnonymOccurrences...),
		analysis.ExecOccurrences...),
		analysis.PluginOccurrences...),
		analysis.GoGenerateOccurrences...),
		analysis.UnsafeOccurrences...),
		analysis.CgoOccurrences...),
		analysis.IndirectOccurrences...),
		analysis.ReflectOccurrences...)
	// Print occurrences
	analysis.PrintOccurrences(analysis.IndirectOccurrences)
	// analysis.PrintOccurrences(occurrences)

	// Count unique occurrences
	initCount, anonymCount, osExecCount, pluginCount, goGenerateCount, unsafeCount, cgoCount, indirectCount, reflectCount := analysis.CountUniqueOccurrences(occurrences)
	fmt.Println()
	fmt.Println()
	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Printf("║ Attack Surface Analysis: %s			       ║\n", filepath.Base(strings.TrimSuffix(modulePath, "/"+filepath.Base(modulePath))))
	fmt.Println("╠══════════════════════════════════════════════════════════════╣")
	fmt.Printf("║ Unique occurrences of init() function:            %10d ║\n", initCount)
	fmt.Printf("║ Initialization with anonymous function:           %10d ║\n", anonymCount)
	fmt.Printf("║ Invocation from the os/exec package:              %10d ║\n", osExecCount)
	fmt.Printf("║ Plugin dynamically loaded:                        %10d ║\n", pluginCount)
	fmt.Printf("║ go:generate directive:                            %10d ║\n", goGenerateCount)
	fmt.Printf("║ Unsafe pointers:                                  %10d ║\n", unsafeCount)
	fmt.Printf("║ CGO pointers:                                     %10d ║\n", cgoCount)
	fmt.Printf("║ Indirect method calls via interfaces:             %10d ║\n", indirectCount)
	fmt.Printf("║ Invocation of reflection:			    %10d ║\n", reflectCount)
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
}

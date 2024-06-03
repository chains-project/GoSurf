package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	analysis "github.com/chains-project/capslock-analysis/gosurface/libs"
)

var (
	initOccurrences        []*analysis.Occurrence
	anonymOccurrences      []*analysis.Occurrence
	execOccurrences        []*analysis.Occurrence
	pluginOccurrences      []*analysis.Occurrence
	goGenerateOccurrences  []*analysis.Occurrence
	goTestOccurrences      []*analysis.Occurrence
	unsafeOccurrences      []*analysis.Occurrence
	cgoOccurrences         []*analysis.Occurrence
	indirectOccurrences    []*analysis.Occurrence
	reflectOccurrences     []*analysis.Occurrence
	constructorOccurrences []*analysis.Occurrence
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
	analysis.AnalyzeModule(modulePath, &initOccurrences, analysis.InitFuncParser{})
	analysis.AnalyzeModule(modulePath, &anonymOccurrences, analysis.AnonymFuncParser{})
	analysis.AnalyzeModule(modulePath, &execOccurrences, analysis.ExecParser{})
	analysis.AnalyzeModule(modulePath, &pluginOccurrences, analysis.PluginParser{})
	analysis.AnalyzeModule(modulePath, &goGenerateOccurrences, analysis.GoGenerateParser{})
	analysis.AnalyzeModule(modulePath, &goTestOccurrences, analysis.GoTestParser{})
	analysis.AnalyzeModule(modulePath, &unsafeOccurrences, analysis.UnsafeParser{})
	analysis.AnalyzeModule(modulePath, &cgoOccurrences, analysis.CgoParser{})
	analysis.AnalyzeModule(modulePath, &indirectOccurrences, analysis.IndirectParser{})
	analysis.AnalyzeModule(modulePath, &reflectOccurrences, analysis.ReflectParser{})
	//analysis.AnalyzeModule(modulePath, &constructorOccurrences, analysis.ConstructorParser{})

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
	occurrences := append(append(append(append(append(append(append(append(append(append(
		initOccurrences,
		anonymOccurrences...),
		execOccurrences...),
		pluginOccurrences...),
		goGenerateOccurrences...),
		goTestOccurrences...),
		unsafeOccurrences...),
		cgoOccurrences...),
		indirectOccurrences...),
		reflectOccurrences...),
		constructorOccurrences...)

	// Print occurrences
	analysis.PrintOccurrences(constructorOccurrences)
	// analysis.PrintOccurrences(occurrences)

	// Count unique occurrences
	initCount, anonymCount, execCount, pluginCount, goGenerateCount, goTestCount, unsafeCount, cgoCount, indirectCount, reflectCount, constructorCount := analysis.CountUniqueOccurrences(occurrences)
	fmt.Println()
	fmt.Println()
	fmt.Println("╔═════════════════════════════════════════════════════════════════════════╗")
	fmt.Printf("║ Attack Surface Analysis: %s	     		         ║\n", filepath.Base(strings.TrimSuffix(modulePath, "/"+filepath.Base(modulePath))))
	fmt.Println("╠═════════════════════════════════════════════════════════════════════════╣")
	fmt.Printf("║ init() function definitions:                                 %10d ║\n", initCount)
	fmt.Printf("║ global var initialization with anonymous functions:          %10d ║\n", anonymCount)
	fmt.Printf("║ exec function invocations:                                   %10d ║\n", execCount)
	fmt.Printf("║ plugin dynamically loaded:                                   %10d ║\n", pluginCount)
	fmt.Printf("║ 'go:generate' directive usage:                               %10d ║\n", goGenerateCount)
	fmt.Printf("║ testing function definitions:                                %10d ║\n", goTestCount)
	fmt.Printf("║ Unsafe pointers:                                             %10d ║\n", unsafeCount) // TODO: define better
	fmt.Printf("║ C function invocations via CGO:                              %10d ║\n", cgoCount)
	fmt.Printf("║ Indirect method calls via interfaces:                        %10d ║\n", indirectCount)
	fmt.Printf("║ Usage of reflection:                                         %10d ║\n", reflectCount) // TODO: define better
	fmt.Printf("║ Invocation of constructors:                                  %10d ║\n", constructorCount)
	fmt.Println("╚═════════════════════════════════════════════════════════════════════════╝")
}

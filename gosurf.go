package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	analysis "example.com/gosurf/gosurf/libs"
)

var (
	initOccurrences        []*analysis.Occurrence
	globalVarOccurrences   []*analysis.Occurrence
	execOccurrences        []*analysis.Occurrence
	pluginOccurrences      []*analysis.Occurrence
	goGenerateOccurrences  []*analysis.Occurrence
	goTestOccurrences      []*analysis.Occurrence
	unsafeOccurrences      []*analysis.Occurrence
	cgoOccurrences         []*analysis.Occurrence
	interfaceOccurrences   []*analysis.Occurrence
	reflectOccurrences     []*analysis.Occurrence
	constructorOccurrences []*analysis.Occurrence
	assemblyOccurrences    []*analysis.Occurrence
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

	fmt.Println("GoSurf is a tool that aims to analyze the potential attack surface of open-source Go packages and modules.")
	fmt.Println("It looks for occurrences of various features and constructs that could potentially introduce security risks.")
	fmt.Println()

	// TODO: currently only fetches direct dependencies in the module, not external dependencies

	// Get direct dependencies
	direct_dependencies, err := analysis.GetDependencies(modulePath)
	if err != nil {
		fmt.Printf("Error getting files in module: %v\n", err)
		return
	}
	// analysis.PrintDependencies(direct_dependencies)

	// Analyze all the module direct dependencies
	for _, dep := range direct_dependencies {
		analysis.AnalyzePackage(dep, &initOccurrences, analysis.InitFuncParser{})
		analysis.AnalyzePackage(dep, &globalVarOccurrences, analysis.GlobalVarParser{})
		analysis.AnalyzePackage(dep, &execOccurrences, analysis.ExecParser{})
		analysis.AnalyzePackage(dep, &pluginOccurrences, analysis.PluginParser{})
		analysis.AnalyzePackage(dep, &goGenerateOccurrences, analysis.GoGenerateParser{})
		analysis.AnalyzePackage(dep, &goTestOccurrences, analysis.GoTestParser{})
		analysis.AnalyzePackage(dep, &unsafeOccurrences, analysis.UnsafeParser{})
		analysis.AnalyzePackage(dep, &cgoOccurrences, analysis.CgoParser{})
		analysis.AnalyzePackage(dep, &interfaceOccurrences, analysis.InterfaceParser{})
		analysis.AnalyzePackage(dep, &reflectOccurrences, analysis.ReflectParser{})
		analysis.AnalyzePackage(dep, &constructorOccurrences, analysis.ConstructorParser{})
		analysis.AnalyzePackage(dep, &assemblyOccurrences, analysis.AssemblyParser{})
	}

	// Convert occurrences to JSON
	occurrences := append(append(append(append(append(append(append(append(append(append(append(
		initOccurrences,
		globalVarOccurrences...),
		execOccurrences...),
		pluginOccurrences...),
		goGenerateOccurrences...),
		goTestOccurrences...),
		unsafeOccurrences...),
		cgoOccurrences...),
		interfaceOccurrences...),
		reflectOccurrences...),
		constructorOccurrences...),
		assemblyOccurrences...)

	// Print occurrences
	// analysis.PrintOccurrences(assemblyOccurrences)
	// analysis.PrintOccurrences(occurrences)

	// Count unique occurrences
	initCount, globalVarCount, execCount, pluginCount, goGenerateCount, goTestCount, unsafeCount, cgoCount, interfaceCount, reflectCount, constructorCount, assemblyCount := analysis.CountUniqueOccurrences(occurrences)
	fmt.Println()
	fmt.Println()
	fmt.Println("╔═════════════════════════════════════════════════════════════════════════╗")
	fmt.Printf("║ Attack Surface Analysis: %s	     		         ║\n", filepath.Base(strings.TrimSuffix(modulePath, "/"+filepath.Base(modulePath))))
	fmt.Println("╠═════════════════════════════════════════════════════════════════════════╣")
	fmt.Printf("║ [P1] Static Code Generation:                                 %10d ║\n", goGenerateCount)
	fmt.Printf("║ [P2] Testing Functions:                                      %10d ║\n", goTestCount)
	fmt.Printf("║ [I1] Global Variable Initialization:                         %10d ║\n", globalVarCount)
	fmt.Printf("║ [I2] init() Functions:                                       %10d ║\n", initCount)
	fmt.Printf("║ [E1] Constructor Methods:                                    %10d ║\n", constructorCount)
	fmt.Printf("║ [E2] Reflection:                                             %10d ║\n", reflectCount)
	fmt.Printf("║ [E3] Interfaces:                                             %10d ║\n", interfaceCount) // TODO: define better
	fmt.Printf("║ [E4] Unsafe Pointers:                                        %10d ║\n", unsafeCount)
	fmt.Printf("║ [E5] CGO Functions:                                          %10d ║\n", cgoCount)
	fmt.Printf("║ [E6] Assembly Functions:                                     %10d ║\n", assemblyCount) // TODO: define better
	fmt.Printf("║ [E7] Dynamic Plugins:                                        %10d ║\n", pluginCount)
	fmt.Printf("║ [E8] External Execution:                                     %10d ║\n", execCount)
	fmt.Println("╚═════════════════════════════════════════════════════════════════════════╝")
}

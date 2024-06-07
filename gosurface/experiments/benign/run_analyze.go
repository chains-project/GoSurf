package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	analysis "github.com/chains-project/capslock-analysis/gosurface/libs"
)

type ModuleDetails struct {
	ModulePath             string
	Version                string
	InitCount              int
	GlobalVarCount         int
	ExecCount              int
	PluginCount            int
	GoGenerateCount        int
	GoTestCount            int
	UnsafeCount            int
	CgoCount               int
	IndirectCount          int
	ReflectCount           int
	ConstructorCount       int
	AssemblyCount          int
	InitOccurrences        []*analysis.Occurrence
	GlobalVarOccurrences   []*analysis.Occurrence
	ExecOccurrences        []*analysis.Occurrence
	PluginOccurrences      []*analysis.Occurrence
	GoGenerateOccurrences  []*analysis.Occurrence
	GoTestOccurrences      []*analysis.Occurrence
	UnsafeOccurrences      []*analysis.Occurrence
	CgoOccurrences         []*analysis.Occurrence
	IndirectOccurrences    []*analysis.Occurrence
	ReflectOccurrences     []*analysis.Occurrence
	ConstructorOccurrences []*analysis.Occurrence
	AssemblyOccurrences    []*analysis.Occurrence
}

func main() {

	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}
	moduleFilePath := filepath.Join(currentDir, "modules_list.json")

	// Modules gathering from libraries.io
	var allModules []map[string]interface{}

	// Retrieve TOP x packages from libraries.io API
	/*
		for page := 1; page <= 5; page++ {

			url := fmt.Sprintf("https://libraries.io/api/search?order=desc&platforms=Go&sort=dependents_count&per_page=1&page=%d&api_key=ff76aa15a1d65e44843fb94dab1ead62", page)
			resp, err := http.Get(url)
			if err != nil {
				fmt.Println("Error making HTTP request:", err)
				return
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Error reading response body:", err)
				return
			}

			var data []map[string]interface{}
			err = json.Unmarshal(body, &data)
			if err != nil {
				fmt.Println("Error unmarshaling JSON:", err)
				return
			}

			// Remove the "versions" field from each map
			for i := range data {
				delete(data[i], "versions")
			}

			allModules = append(allModules, data...)
		}

		// Print retrieved packages to JSON
		modifiedJSON, err := json.MarshalIndent(allModules, "", "  ")
		if err != nil {
			fmt.Println("Error marshaling JSON:", err)
			return
		}
		err = os.WriteFile(moduleFilePath, modifiedJSON, 0644)
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
		fmt.Println("Modified JSON written to", moduleFilePath)
		itemCount := len(allModules)
		fmt.Println("Number of items in the JSON file:", itemCount)
	*/

	/* Read package paths from file*/
	file, err := os.Open(moduleFilePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&allModules)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}
	itemCount := len(allModules)
	/* **********************/

	// Parse the HTML templates
	overviewTmpl, err := template.ParseFiles("../report_tmpl/tmpl_overview.html")
	if err != nil {
		fmt.Println("Error parsing overview template:", err)
		return
	}

	detailsTmpl, err := template.ParseFiles("../report_tmpl/tmpl_details.html")
	if err != nil {
		fmt.Println("Error parsing details template:", err)
		return
	}

	// Create the HTML files
	overviewFile, err := os.Create("results_overview.html")
	if err != nil {
		fmt.Println("Error creating overview file:", err)
		return
	}
	defer overviewFile.Close()

	detailsFile, err := os.Create("results_detail.html")
	if err != nil {
		fmt.Println("Error creating details file:", err)
		return
	}
	defer detailsFile.Close()

	// Slices to hold PackageAnalysis and PackageDetails instances
	var moduleDetails []ModuleDetails

	// Get and analyze each module
	for i, pkg := range allModules {
		packageManagerURL := pkg["package_manager_url"].(string)
		latestReleaseNumber := pkg["latest_release_number"].(string)

		// Construct the package import path and version
		importPath := strings.TrimPrefix(packageManagerURL, "https://pkg.go.dev/")
		version := "@" + latestReleaseNumber
		fmt.Printf("\n[%d/%d] Analyzing package %s...\n", i+1, itemCount, importPath)

		// Get the module
		cmd := exec.Command("go", "get", importPath+version)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			fmt.Printf("Error getting package %s: %v\n", importPath, err)
			continue
		}

		// Analyze the module
		var initOccurrences, globalVarOccurrences, execOccurrences, pluginOccurrences, goGenerateOccurrences, goTestOccurrences, unsafeOccurrences, cgoOccurrences, indirectOccurrences, reflectOccurrences, constructorOccurrences, assemblyOccurrences []*analysis.Occurrence
		modulePath := filepath.Join(os.Getenv("GOPATH"), "pkg/mod", importPath+"@"+latestReleaseNumber)

		// Analyze the module and its direct dependencies
		direct_dependencies, err := analysis.GetDependencies(modulePath)
		if err != nil {
			fmt.Printf("Error getting files in module: %v\n", err)
			return
		}

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
			analysis.AnalyzePackage(dep, &indirectOccurrences, analysis.IndirectParser{})
			analysis.AnalyzePackage(dep, &reflectOccurrences, analysis.ReflectParser{})
			//analysis.AnalyzePackage(dep, &constructorOccurrences, analysis.ConstructorParser{})
		}

		// Convert occurrences to JSON
		occurrences := append(append(append(append(append(append(append(append(append(append(
			initOccurrences,
			globalVarOccurrences...),
			execOccurrences...),
			pluginOccurrences...),
			goGenerateOccurrences...),
			goTestOccurrences...),
			unsafeOccurrences...),
			cgoOccurrences...),
			indirectOccurrences...),
			reflectOccurrences...),
			constructorOccurrences...)

		// Count unique occurrences
		initCount, globalVarCount, execCount, pluginCount, goGenerateCount, goTestCount, unsafeCount, cgoCount, indirectCount, reflectCount, constructorCount, assemblyCount := analysis.CountUniqueOccurrences(occurrences)

		// Create a ModuleDetails instance and append it to the slice
		moduleDetail := ModuleDetails{
			ModulePath:             modulePath,
			Version:                latestReleaseNumber,
			InitCount:              initCount,
			GlobalVarCount:         globalVarCount,
			ExecCount:              execCount,
			PluginCount:            pluginCount,
			GoGenerateCount:        goGenerateCount,
			GoTestCount:            goTestCount,
			UnsafeCount:            unsafeCount,
			CgoCount:               cgoCount,
			IndirectCount:          indirectCount,
			ReflectCount:           reflectCount,
			ConstructorCount:       constructorCount,
			AssemblyCount:          assemblyCount,
			InitOccurrences:        initOccurrences,
			GlobalVarOccurrences:   globalVarOccurrences,
			ExecOccurrences:        execOccurrences,
			PluginOccurrences:      pluginOccurrences,
			GoGenerateOccurrences:  goGenerateOccurrences,
			GoTestOccurrences:      goTestOccurrences,
			UnsafeOccurrences:      unsafeOccurrences,
			CgoOccurrences:         cgoOccurrences,
			IndirectOccurrences:    indirectOccurrences,
			ReflectOccurrences:     reflectOccurrences,
			ConstructorOccurrences: constructorOccurrences,
			AssemblyOccurrences:    assemblyOccurrences,
		}
		moduleDetails = append(moduleDetails, moduleDetail)
	}

	// Execute the overview template with the ModuleDetails instances
	err = overviewTmpl.Execute(overviewFile, moduleDetails)
	if err != nil {
		fmt.Println("Error executing overview template:", err)
		return
	}

	// Execute the details template with the ModuleDetails instances
	err = detailsTmpl.Execute(detailsFile, moduleDetails)
	if err != nil {
		fmt.Println("Error executing details template:", err)
		return
	}

	fmt.Println("HTML report generated successfully in the current directory.")

}

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

type PackageDetails struct {
	PackagePath            string
	Version                string
	InitCount              int
	AnonymCount            int
	ExecCount              int
	PluginCount            int
	GoGenerateCount        int
	GoTestCount            int
	UnsafeCount            int
	CgoCount               int
	IndirectCount          int
	ReflectCount           int
	ConstructorCount       int
	InitOccurrences        []*analysis.Occurrence
	AnonymOccurrences      []*analysis.Occurrence
	ExecOccurrences        []*analysis.Occurrence
	PluginOccurrences      []*analysis.Occurrence
	GoGenerateOccurrences  []*analysis.Occurrence
	GoTestOccurrences      []*analysis.Occurrence
	UnsafeOccurrences      []*analysis.Occurrence
	CgoOccurrences         []*analysis.Occurrence
	IndirectOccurrences    []*analysis.Occurrence
	ReflectOccurrences     []*analysis.Occurrence
	ConstructorOccurrences []*analysis.Occurrence
}

func main() {

	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}
	filePath := filepath.Join(currentDir, "results.json")

	/*
		var allPackages []map[string]interface{}

		// Retrieve TOP x packages from libraries.io API
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

			allPackages = append(allPackages, data...)
		}

		// Print retrieved packages to JSON
		modifiedJSON, err := json.MarshalIndent(allPackages, "", "  ")
		if err != nil {
			fmt.Println("Error marshaling JSON:", err)
			return
		}
		err = os.WriteFile(filePath, modifiedJSON, 0644)
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
		fmt.Println("Modified JSON written to", filePath)
		itemCount := len(allPackages)
		fmt.Println("Number of items in the JSON file:", itemCount)

	*/

	/* Start Testing code*/

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	var allPackages []map[string]interface{}
	err = json.NewDecoder(file).Decode(&allPackages)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}
	itemCount := len(allPackages)
	/* End Testing code*/

	// Parse the HTML templates
	overviewTmpl, err := template.ParseFiles("tmpl_overview.html")
	if err != nil {
		fmt.Println("Error parsing overview template:", err)
		return
	}

	detailsTmpl, err := template.ParseFiles("tmpl_details.html")
	if err != nil {
		fmt.Println("Error parsing details template:", err)
		return
	}

	// Create the HTML files
	overviewFile, err := os.Create("analysis_overview.html")
	if err != nil {
		fmt.Println("Error creating overview file:", err)
		return
	}
	defer overviewFile.Close()

	detailsFile, err := os.Create("analysis_detail.html")
	if err != nil {
		fmt.Println("Error creating details file:", err)
		return
	}
	defer detailsFile.Close()

	// Slices to hold PackageAnalysis and PackageDetails instances
	var packageDetails []PackageDetails

	// Get and analyze each package
	for i, pkg := range allPackages {
		packageManagerURL := pkg["package_manager_url"].(string)
		latestReleaseNumber := pkg["latest_release_number"].(string)

		// Construct the package import path and version
		importPath := strings.TrimPrefix(packageManagerURL, "https://pkg.go.dev/")
		version := "@" + latestReleaseNumber
		fmt.Printf("[%d/%d] Analyzing package %s...\n", i+1, itemCount, importPath)

		// Get the package
		cmd := exec.Command("go", "get", importPath+version)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			fmt.Printf("Error getting package %s: %v\n", importPath, err)
			continue
		}

		// Analyze the package
		var initOccurrences, anonymOccurrences, execOccurrences, pluginOccurrences, goGenerateOccurrences, goTestOccurrences, unsafeOccurrences, cgoOccurrences, indirectOccurrences, reflectOccurrences, constructorOccurrences []*analysis.Occurrence
		packagePath := filepath.Join(os.Getenv("GOPATH"), "pkg/mod", importPath+"@"+latestReleaseNumber)

		// Analyze the module and its direct dependencies
		analysis.AnalyzeModule(packagePath, &initOccurrences, analysis.InitFuncParser{})
		analysis.AnalyzeModule(packagePath, &anonymOccurrences, analysis.AnonymFuncParser{})
		analysis.AnalyzeModule(packagePath, &execOccurrences, analysis.ExecParser{})
		analysis.AnalyzeModule(packagePath, &pluginOccurrences, analysis.PluginParser{})
		analysis.AnalyzeModule(packagePath, &goGenerateOccurrences, analysis.GoGenerateParser{})
		analysis.AnalyzeModule(packagePath, &goTestOccurrences, analysis.GoTestParser{})
		analysis.AnalyzeModule(packagePath, &unsafeOccurrences, analysis.UnsafeParser{})
		analysis.AnalyzeModule(packagePath, &cgoOccurrences, analysis.CgoParser{})
		analysis.AnalyzeModule(packagePath, &indirectOccurrences, analysis.IndirectParser{})
		analysis.AnalyzeModule(packagePath, &reflectOccurrences, analysis.ReflectParser{})
		// analysis.AnalyzeModule(packagePath, &constructorOccurrences, analysis.ConstructorParser{})

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

		// Count unique occurrences
		initCount, anonymCount, execCount, pluginCount, goGenerateCount, goTestCount, unsafeCount, cgoCount, indirectCount, reflectCount, constructorCount := analysis.CountUniqueOccurrences(occurrences)

		// Create a PackageDetails instance and append it to the slice
		packageDetail := PackageDetails{
			PackagePath:            packagePath,
			Version:                latestReleaseNumber,
			InitCount:              initCount,
			AnonymCount:            anonymCount,
			ExecCount:              execCount,
			PluginCount:            pluginCount,
			GoGenerateCount:        goGenerateCount,
			GoTestCount:            goTestCount,
			UnsafeCount:            unsafeCount,
			CgoCount:               cgoCount,
			IndirectCount:          indirectCount,
			ReflectCount:           reflectCount,
			ConstructorCount:       constructorCount,
			InitOccurrences:        initOccurrences,
			AnonymOccurrences:      anonymOccurrences,
			ExecOccurrences:        execOccurrences,
			PluginOccurrences:      pluginOccurrences,
			GoGenerateOccurrences:  goGenerateOccurrences,
			GoTestOccurrences:      goTestOccurrences,
			UnsafeOccurrences:      unsafeOccurrences,
			CgoOccurrences:         cgoOccurrences,
			IndirectOccurrences:    indirectOccurrences,
			ReflectOccurrences:     reflectOccurrences,
			ConstructorOccurrences: constructorOccurrences,
		}
		packageDetails = append(packageDetails, packageDetail)
	}

	// Execute the overview template with the PackageAnalysis instances
	err = overviewTmpl.Execute(overviewFile, packageDetails)
	if err != nil {
		fmt.Println("Error executing overview template:", err)
		return
	}

	// Execute the details template with the PackageDetails instances
	err = detailsTmpl.Execute(detailsFile, packageDetails)
	if err != nil {
		fmt.Println("Error executing details template:", err)
		return
	}

	fmt.Println("HTML file generated: attack_surface_analysis.html")

}

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	analysis "github.com/chains-project/capslock-analysis/gosurface/libs"
)

type Repository struct {
	ContributionsCount             int      `json:"contributions_count"`
	DependentReposCount            int      `json:"dependent_repos_count"`
	DependentsCount                int      `json:"dependents_count"`
	Description                    string   `json:"description"`
	Forks                          int      `json:"forks"`
	Keywords                       []string `json:"keywords"`
	Language                       string   `json:"language"`
	LatestReleaseNumber            string   `json:"latest_release_number"`
	LatestReleasePublishedAt       string   `json:"latest_release_published_at"`
	LatestStableReleaseNumber      string   `json:"latest_stable_release_number"`
	LatestStableReleasePublishedAt string   `json:"latest_stable_release_published_at"`
	Name                           string   `json:"name"`
	PackageManagerURL              string   `json:"package_manager_url"`
	Platform                       string   `json:"platform"`
	RepositoryURL                  string   `json:"repository_url"`
	Stars                          int      `json:"stars"`
}

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
}

func main() {

	// Create folders
	if err := os.MkdirAll("./results", 0755); err != nil {
		fmt.Printf("Error creating results directory: %v\n", err)
		return
	}

	// Retrieve module information from libraries.io and write to modules_info.json
	var moduleInfoFile = "./results/modules_info.json"
	retrieveModulesFromLibrariesIO(moduleInfoFile)

	// Read package paths from modules_info.json
	allModules, itemCount, err := readModulesFromFile(moduleInfoFile)
	if err != nil {
		fmt.Printf("Error reading modules from file: %v\n", err)
		return
	}

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
	overviewFile, err := os.Create("./results/results_overview.html")
	if err != nil {
		fmt.Println("Error creating overview file:", err)
		return
	}
	defer overviewFile.Close()
	detailsFile, err := os.Create("./results/results_detail.html")
	if err != nil {
		fmt.Println("Error creating details file:", err)
		return
	}
	defer detailsFile.Close()

	// Get each module
	for idx, module := range allModules {
		packageManagerURL := module.PackageManagerURL
		latestReleaseNumber := module.LatestReleaseNumber

		// Construct the module import path and version
		importPath := strings.TrimPrefix(packageManagerURL, "https://pkg.go.dev/")
		version := "@" + latestReleaseNumber
		fmt.Printf("\n[%d/%d] Getting module %s...\n", idx+1, itemCount, importPath)

		// Get the module
		cmd := exec.Command("go", "get", importPath+version)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			fmt.Printf("Error getting module %s: %v\n", importPath, err)
			continue
		}
	}

	var moduleDetailsList []ModuleDetails

	// Analyze each module
	for i, module := range allModules {

		packageManagerURL := module.PackageManagerURL
		latestReleaseNumber := module.LatestReleaseNumber

		// Construct the module import path and version
		importPath := strings.TrimPrefix(packageManagerURL, "https://pkg.go.dev/")
		fmt.Printf("\n[%d/%d] Analyzing module %s...\n", i+1, itemCount, importPath)

		// Analyze the module
		var initOccurrences, globalVarOccurrences, execOccurrences, pluginOccurrences, goGenerateOccurrences, goTestOccurrences, unsafeOccurrences, cgoOccurrences, indirectOccurrences, reflectOccurrences, constructorOccurrences []*analysis.Occurrence
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
		initCount, globalVarCount, execCount, pluginCount, goGenerateCount, goTestCount, unsafeCount, cgoCount, indirectCount, reflectCount, constructorCount := analysis.CountUniqueOccurrences(occurrences)

		// Create a ModuleDetails instance and append it to the slice
		moduleDetails := ModuleDetails{
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
		}
		moduleDetailsList = append(moduleDetailsList, moduleDetails)
	}

	// Execute the overview template with the ModuleDetails instances
	err = overviewTmpl.Execute(overviewFile, moduleDetailsList)
	if err != nil {
		fmt.Println("Error executing overview template:", err)
		return
	}

	// Execute the details template with the ModuleDetails instances
	err = detailsTmpl.Execute(detailsFile, moduleDetailsList)
	if err != nil {
		fmt.Println("Error executing details template:", err)
		return
	}

	fmt.Println("\nHTML report generated successfully in the current directory.")

}

func retrieveModulesFromLibrariesIO(moduleInfoFile string) {

	fmt.Printf("Retrieving of TOP x modules from libraries.io API\n")

	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}
	moduleInfoFile = filepath.Join(currentDir, moduleInfoFile)

	var allModules []Repository

	// Retrieve TOP x packages from libraries.io API
	for page := 1; page <= 1; page++ {
		url := fmt.Sprintf("https://libraries.io/api/search?order=desc&platforms=Go&sort=dependents_count&per_page=5&page=%d&api_key=ff76aa15a1d65e44843fb94dab1ead62", page)
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

		var data []Repository
		err = json.Unmarshal(body, &data)
		if err != nil {
			fmt.Println("Error unmarshaling JSON:", err)
			return
		}

		allModules = append(allModules, data...)
	}

	// Write retrieved packages to modules_list.json
	jsonData, err := json.MarshalIndent(allModules, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}
	err = os.WriteFile(moduleInfoFile, jsonData, 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
	fmt.Println("Retrieved information written to JSON file", moduleInfoFile)

}

func readModulesFromFile(moduleFilePath string) ([]Repository, int, error) {
	var allModules []Repository

	file, err := os.Open(moduleFilePath)
	if err != nil {
		return nil, 0, err
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&allModules)
	if err != nil {
		return nil, 0, err
	}

	itemCount := len(allModules)

	return allModules, itemCount, nil
}

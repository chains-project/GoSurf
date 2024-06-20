package main

import (
	"bufio"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	analysis "github.com/chains-project/capslock-analysis/gosurface/libs"
)

type ModuleDetails struct {
	Name                   string
	ModulePath             string
	RepositoryURL          string
	Dependants             string
	Version                string
	LOC                    int
	InitCount              []float64
	GlobalVarCount         []float64
	ExecCount              []float64
	PluginCount            []float64
	GoGenerateCount        []float64
	GoTestCount            []float64
	UnsafeCount            []float64
	CgoCount               []float64
	InterfaceCount         []float64
	ReflectCount           []float64
	ConstructorCount       []float64
	AssemblyCount          []float64
	InitOccurrences        []*analysis.Occurrence
	GlobalVarOccurrences   []*analysis.Occurrence
	ExecOccurrences        []*analysis.Occurrence
	PluginOccurrences      []*analysis.Occurrence
	GoGenerateOccurrences  []*analysis.Occurrence
	GoTestOccurrences      []*analysis.Occurrence
	UnsafeOccurrences      []*analysis.Occurrence
	CgoOccurrences         []*analysis.Occurrence
	InterfaceOccurrences   []*analysis.Occurrence
	ReflectOccurrences     []*analysis.Occurrence
	ConstructorOccurrences []*analysis.Occurrence
	AssemblyOccurrences    []*analysis.Occurrence
}

func main() {

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("GITHUB_TOKEN environment variable is not set")
	}

	if err := os.MkdirAll("./modules", 0755); err != nil {
		fmt.Printf("Error creating cloned repos directory: %v\n", err)
		return
	}
	if err := os.MkdirAll("./results", 0755); err != nil {
		fmt.Printf("Error creating results directory: %v\n", err)
		return
	}

	var urlListFile = "repo_urls.txt"
	allModules, itemCount, err := readModulesFromFile(urlListFile)
	if err != nil {
		fmt.Printf("Error reading modules from file: %v\n", err)
		return
	}

	overviewTmpl, err := template.ParseFiles("./web/tmpl_overview.html")
	if err != nil {
		fmt.Println("Error parsing overview template:", err)
		return
	}
	detailsTmpl, err := template.ParseFiles("./web/tmpl_details.html")
	if err != nil {
		fmt.Println("Error parsing details template:", err)
		return
	}
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

	// Clone each module from github and run the analysis
	for _, module := range allModules {

		repoName := module.Name
		repoURL := module.RepositoryURL
		releaseNumber := module.Version
		modulePath := filepath.Join("./modules", fmt.Sprintf("%s@%s", repoName, releaseNumber))

		fmt.Printf("\n\nCloning module %s@%s into %s...\n", repoName, releaseNumber, modulePath)
		if _, err := os.Stat(modulePath); !os.IsNotExist(err) {
			fmt.Printf("Skipping %s as %s already exists.\n", repoURL, modulePath)
			continue
		}

		cmd := exec.Command("git", "clone", "--branch", releaseNumber, repoURL, modulePath)
		fmt.Printf("Executing command: %s\n", cmd)
		err := cmd.Run()
		if err != nil {
			fmt.Println("Error cloning repository:", err)
			return
		}
	}

	// Run the analysis
	for idx, module := range allModules {

		repoName := module.Name
		releaseNumber := module.Version
		modulePath := filepath.Join("./modules", fmt.Sprintf("%s@%s", repoName, releaseNumber))

		fmt.Printf("\n[%d/%d] Analyzing module %s@%s...\n", idx+1, itemCount, repoName, releaseNumber)

		// TODO: use directly the API of this package
		locCount, err := analysis.GetLineOfCodeCount(modulePath)
		if err != nil {
			fmt.Printf("Error getting line of code count for %s: %v\n", repoName, err)
			continue
		}

		// Analyze the module and its direct dependencies
		var initOccurrences, globalVarOccurrences, execOccurrences, pluginOccurrences, goGenerateOccurrences, goTestOccurrences, unsafeOccurrences, cgoOccurrences, interfaceOccurrences, reflectOccurrences, constructorOccurrences, assemblyOccurrences []*analysis.Occurrence
		direct_dependencies, err := analysis.GetDependencies(modulePath)
		if err != nil {
			fmt.Printf("Error getting files in module: %v\n", err)
			return
		}

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

		initCount, globalVarCount, execCount, pluginCount, goGenerateCount, goTestCount, unsafeCount, cgoCount, interfaceCount, reflectCount, constructorCount, assemblyCount := analysis.CountUniqueOccurrences(occurrences)

		moduleDetails := ModuleDetails{
			ModulePath:             modulePath,
			Version:                releaseNumber,
			LOC:                    locCount,
			InitCount:              []float64{float64(initCount), float64(initCount) / float64(locCount)},
			GlobalVarCount:         []float64{float64(globalVarCount), float64(globalVarCount) / float64(locCount)},
			ExecCount:              []float64{float64(execCount), float64(execCount) / float64(locCount)},
			PluginCount:            []float64{float64(pluginCount), float64(pluginCount) / float64(locCount)},
			GoGenerateCount:        []float64{float64(goGenerateCount), float64(goGenerateCount) / float64(locCount)},
			GoTestCount:            []float64{float64(goTestCount), float64(goTestCount) / float64(locCount)},
			UnsafeCount:            []float64{float64(unsafeCount), float64(unsafeCount) / float64(locCount)},
			CgoCount:               []float64{float64(cgoCount), float64(cgoCount) / float64(locCount)},
			InterfaceCount:         []float64{float64(interfaceCount), float64(interfaceCount) / float64(locCount)},
			ReflectCount:           []float64{float64(reflectCount), float64(reflectCount) / float64(locCount)},
			ConstructorCount:       []float64{float64(constructorCount), float64(constructorCount) / float64(locCount)},
			AssemblyCount:          []float64{float64(assemblyCount), float64(assemblyCount) / float64(locCount)},
			InitOccurrences:        initOccurrences,
			GlobalVarOccurrences:   globalVarOccurrences,
			ExecOccurrences:        execOccurrences,
			PluginOccurrences:      pluginOccurrences,
			GoGenerateOccurrences:  goGenerateOccurrences,
			GoTestOccurrences:      goTestOccurrences,
			UnsafeOccurrences:      unsafeOccurrences,
			CgoOccurrences:         cgoOccurrences,
			InterfaceOccurrences:   interfaceOccurrences,
			ReflectOccurrences:     reflectOccurrences,
			ConstructorOccurrences: constructorOccurrences,
			AssemblyOccurrences:    assemblyOccurrences,
		}
		allModules[idx] = moduleDetails
	}

	// Execute the template with the ModuleDetails instances
	err = overviewTmpl.Execute(overviewFile, allModules)
	if err != nil {
		fmt.Println("Error executing overview template:", err)
		return
	}
	err = detailsTmpl.Execute(detailsFile, allModules)
	if err != nil {
		fmt.Println("Error executing details template:", err)
		return
	}

	fmt.Println("\nHTML report generated successfully in the ./results directory.")

}

func readModulesFromFile(urlListFile string) ([]ModuleDetails, int, error) {
	file, err := os.Open(urlListFile)
	if err != nil {
		return nil, 0, err
	}
	defer file.Close()

	var allModules []ModuleDetails
	var itemCount int

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, " ")
		if len(parts) == 2 {
			fullRepoURL, version := parts[0], parts[1]
			parsedURL, err := url.Parse(fullRepoURL)
			if err != nil {
				continue
			}
			repoName := path.Base(parsedURL.Path)
			moduleDetails := ModuleDetails{
				Name:          repoName,
				RepositoryURL: fullRepoURL,
				Version:       version,
			}
			allModules = append(allModules, moduleDetails)
			itemCount++

		}
	}

	if err := scanner.Err(); err != nil {
		return nil, 0, err
	}

	return allModules, itemCount, nil
}

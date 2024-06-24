package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
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
	Dependants             int
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

	expName := "exp1"
	if len(os.Args) > 1 {
		expName = os.Args[1]
	}
	if expName != "exp1" && expName != "exp2" {
		fmt.Println("Invalid input. Please provide 'exp1' or 'exp2'.")
		return
	}

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("GITHUB_TOKEN environment variable is not set")
	}

	if err := os.MkdirAll(fmt.Sprintf("./modules/%s", expName), 0755); err != nil {
		fmt.Printf("Error creating cloned repos directory: %v\n", err)
		return
	}
	if err := os.MkdirAll(fmt.Sprintf("./results/%s", expName), 0755); err != nil {
		fmt.Printf("Error creating results directory: %v\n", err)
		return
	}

	var urlListFile = fmt.Sprintf("urls_%s.txt", expName)
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
	overviewFile, err := os.Create(fmt.Sprintf("./results/%s/results_overview.html", expName))
	if err != nil {
		fmt.Println("Error creating overview file:", err)
		return
	}
	defer overviewFile.Close()
	detailsFile, err := os.Create(fmt.Sprintf("./results/%s/results_detail.html", expName))
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
		modulePath := filepath.Join(fmt.Sprintf("./modules/%s", expName), fmt.Sprintf("%s@%s", repoName, releaseNumber))

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
		repoURL := module.RepositoryURL
		releaseNumber := module.Version
		modulePath := filepath.Join(fmt.Sprintf("./modules/%s", expName), fmt.Sprintf("%s@%s", repoName, releaseNumber))

		fmt.Printf("\n[%d/%d] Analyzing module %s@%s...\n", idx+1, itemCount, repoName, releaseNumber)

		dependantsCount, err := fetchDependantsCount(repoURL)
		if err != nil {
			fmt.Printf("Error fetching dependants count for %s: %v\n", repoName, err)
			continue
		}

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
			Dependants:             dependantsCount,
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

	fmt.Printf("\nHTML report generated successfully in the ./results/%s directory.\n", expName)

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

func fetchDependantsCount(packageURL string) (int, error) {

	packageURL = strings.TrimPrefix(packageURL, "https://")
	packageURL = strings.ReplaceAll(packageURL, "/", "%2F")
	apiURL := fmt.Sprintf("https://libraries.io/api/go/%s?api_key=ff76aa15a1d65e44843fb94dab1ead62", packageURL)

	resp, err := http.Get(apiURL)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return 0, err
	}

	dependentsCount, ok := data["dependents_count"].(float64)
	if !ok {
		return 0, fmt.Errorf("failed to parse dependents_count")
	}

	return int(dependentsCount), nil
}

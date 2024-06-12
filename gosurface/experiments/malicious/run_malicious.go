package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/go-git/go-git/v5"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"

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

	// Read the GitHub token from the environment variable
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("GITHUB_TOKEN environment variable is not set")
	}

	// Create folders
	if err := os.MkdirAll("./cloned_repos", 0755); err != nil {
		fmt.Printf("Error creating cloned repos directory: %v\n", err)
		return
	}
	if err := os.MkdirAll("./results", 0755); err != nil {
		fmt.Printf("Error creating results directory: %v\n", err)
		return
	}

	// Retrieve modules information from github repositories and write to modules_info.json
	var urlListFile = "repo_urls.txt"
	var moduleInfoFile = "./results/modules_info.json"
	_, err := os.Stat(moduleInfoFile)
	if os.IsNotExist(err) {
		retrieveModulesFromGithub(urlListFile, moduleInfoFile, token)
	} else {
		fmt.Printf("Module information file %s already exists, skipping retrieval from GitHub.\n", moduleInfoFile)
	}

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

	// Clone each module from github
	for idx, module := range allModules {
		repoName := module.Name
		repoURL := module.RepositoryURL
		cloneDir := filepath.Join("./cloned_repos", fmt.Sprintf("%02d_%s", idx+1, repoName))

		// Clone the repository
		fmt.Printf("Cloning module %s into %s...\n", repoName, cloneDir)
		if _, err := os.Stat(cloneDir); !os.IsNotExist(err) {
			fmt.Printf("Skipping %s as %s already exists.\n", repoURL, cloneDir)
			continue
		}
		_, err := git.PlainClone(cloneDir, false, &git.CloneOptions{
			URL: repoURL,
		})
		if err != nil {
			fmt.Printf("Error Plain Clone %v\n", err)
			continue
		}
	}

	var moduleDetailsList []ModuleDetails

	// Analyze each module
	for idx, module := range allModules {

		// Construct the module path
		fmt.Printf("\n[%d/%d] Analyzing module %s...\n", idx+1, itemCount, module.Name)
		currentDir, _ := os.Getwd()
		latestReleaseNumber := module.LatestReleaseNumber
		modulePath := filepath.Join(currentDir, "cloned_repos", fmt.Sprintf("%02d_%s", idx+1, module.Name))

		// TODO: use directly the API of this package
		// Get the lines of code count
		locCount, err := analysis.GetLineOfCodeCount(modulePath)
		if err != nil {
			fmt.Printf("Error getting line of code count for %s: %v\n", module.Name, err)
			continue
		}

		// Analyze the module and its direct dependencies
		var initOccurrences, globalVarOccurrences, execOccurrences, pluginOccurrences, goGenerateOccurrences, goTestOccurrences, unsafeOccurrences, cgoOccurrences, interfaceOccurrences, reflectOccurrences, constructorOccurrences, assemblyOccurrences []*analysis.Occurrence
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
			analysis.AnalyzePackage(dep, &interfaceOccurrences, analysis.InterfaceParser{})
			analysis.AnalyzePackage(dep, &reflectOccurrences, analysis.ReflectParser{})
			//analysis.AnalyzePackage(dep, &constructorOccurrences, analysis.ConstructorParser{})
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

		// Count unique occurrences
		initCount, globalVarCount, execCount, pluginCount, goGenerateCount, goTestCount, unsafeCount, cgoCount, interfaceCount, reflectCount, constructorCount, assemblyCount := analysis.CountUniqueOccurrences(occurrences)

		// Create a ModuleDetails instance and append it to the slice
		moduleDetails := ModuleDetails{
			ModulePath:             modulePath,
			Version:                latestReleaseNumber,
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

func retrieveModulesFromGithub(urlListFile string, moduleInfoFile string, token string) {

	fmt.Printf("Retrieving module information from GitHub...\n")

	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}
	moduleInfoFile = filepath.Join(currentDir, moduleInfoFile)

	// Create a new GitHub client with the provided token
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// Open the file containing the repository links
	file, err := os.Open(urlListFile)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	// Read the repository links from the file
	scanner := bufio.NewScanner(file)
	var repositories []Repository
	for scanner.Scan() {
		repoURL := scanner.Text()
		repo, err := getRepositoryInfo(client, repoURL)
		if err != nil {
			log.Printf("Failed to retrieve information for %s: %v", repoURL, err)
			continue
		}
		repositories = append(repositories, *repo)
	}

	// Marshal the repositories to JSON
	jsonData, err := json.MarshalIndent(repositories, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}
	err = os.WriteFile(moduleInfoFile, jsonData, 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Println("Retrieved information from GitHub written to JSON file", moduleInfoFile)
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

func getRepositoryInfo(client *github.Client, repoURL string) (*Repository, error) {
	// Parse the repository owner and name from the URL
	parts := strings.Split(repoURL, "/")
	owner := parts[len(parts)-2]
	repo := parts[len(parts)-1]

	// Get the repository information from the GitHub API
	repository, _, err := client.Repositories.Get(context.Background(), owner, repo)
	if err != nil {
		return nil, err
	}

	// Get the list of contributors from the GitHub API
	contributors, _, err := client.Repositories.ListContributors(context.Background(), owner, repo, nil)
	if err != nil {
		return nil, err
	}
	contributorsCount := len(contributors)

	// Get language information
	var language string
	if repository.Language != nil {
		language = *repository.Language
	}

	// Get the latest release information from the GitHub API
	releases, _, err := client.Repositories.ListReleases(context.Background(), owner, repo, nil)
	if err != nil {
		return nil, err
	}

	var latestRelease, latestStableRelease *github.RepositoryRelease
	for _, release := range releases {
		if release.GetPrerelease() {
			continue
		}
		if latestStableRelease == nil || release.GetPublishedAt().After(latestStableRelease.GetPublishedAt().Time) {
			latestStableRelease = release
		}
		if latestRelease == nil || release.GetPublishedAt().After(latestRelease.GetPublishedAt().Time) {
			latestRelease = release
		}
	}

	// Create the Repository struct with the retrieved information
	repoInfo := &Repository{
		ContributionsCount:             contributorsCount,
		DependentReposCount:            0, // Not available from the GitHub API
		DependentsCount:                int(repository.GetSubscribersCount()),
		Description:                    repository.GetDescription(),
		Forks:                          repository.GetForksCount(),
		Keywords:                       repository.Topics,
		Language:                       language,
		LatestReleaseNumber:            getLatestReleaseNumber(latestRelease),
		LatestReleasePublishedAt:       getLatestReleasePublishedAt(latestRelease),
		LatestStableReleaseNumber:      getLatestReleaseNumber(latestStableRelease),
		LatestStableReleasePublishedAt: getLatestReleasePublishedAt(latestStableRelease),
		Name:                           repository.GetName(),
		PackageManagerURL:              "", // Not available from the GitHub API
		RepositoryURL:                  repository.GetHTMLURL(),
		Stars:                          repository.GetStargazersCount(),
	}

	fmt.Printf("Retrieved repository information for %s\n", repoInfo.Name)

	return repoInfo, nil
}

func getLatestReleaseNumber(release *github.RepositoryRelease) string {
	if release == nil {
		return ""
	}
	return release.GetTagName()
}

func getLatestReleasePublishedAt(release *github.RepositoryRelease) string {
	if release == nil {
		return ""
	}
	return release.GetPublishedAt().Format("2006-01-02T15:04:05.000Z")
}

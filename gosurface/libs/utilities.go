package analysis

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Occurrence struct {
	AttackVector  string
	FilePath      string
	LineNumber    int
	VariableName  string // for anonymous functions
	Command       string // for go:generate directive
	MethodInvoked string // for indirect, exec, plugin, cgo
	TypePassed    string // for indirect
	Pattern       string // for constructors
}

type Dependency struct {
	Name string
	Path string
}

type OccurrenceParser interface {
	FindOccurrences(path string, occurrences *[]*Occurrence)
}

// Gets all go files in given path.
func GetDependencies(modulePath string) ([]Dependency, error) { // TODO should rename this one. If getting dependencies, we look at the go.mod file.

	var dependencies []Dependency

	// Check if the parent folder is a package
	isPackage, packageName, packagePath := isGoPackage(modulePath)
	if isPackage {
		//canBuild, _ := canBuildGoPackage(modulePath)
		//if canBuild {
		dependency := Dependency{Name: packageName, Path: packagePath}
		dependencies = append(dependencies, dependency)
		//}
	}

	// Gather subdirectories
	var subdirs []string
	err := filepath.Walk(modulePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && path != modulePath {
			subdirs = append(subdirs, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	totalSubdirs := len(subdirs)
	processedSubdirs := 0

	// Process each subdirectory
	for _, dirPath := range subdirs {
		isPackage, packageName, packagePath := isGoPackage(dirPath)
		if isPackage {
			//canBuild, _ := canBuildGoPackage(dirPath)
			//if canBuild {
			dependency := Dependency{Name: packageName, Path: packagePath}
			dependencies = append(dependencies, dependency)
			//}
		}
		processedSubdirs++
		updateProgressBar(processedSubdirs, totalSubdirs)
	}

	return dependencies, nil
}

func isGoPackage(dirPath string) (bool, string, string) {
	goFiles := findGoFiles(dirPath)
	for _, goFile := range goFiles {
		filePath := filepath.Join(dirPath, goFile)
		content, err := os.ReadFile(filePath)
		if err != nil {
			continue
		}
		match := packageRegex.FindSubmatch(content)
		if len(match) > 1 {
			packageName := string(match[1])
			return true, packageName, dirPath
		}
	}
	return false, "", ""
}

/*
func canBuildGoPackage(dirPath string) (bool, string) {
	cmd := exec.Command("go", "build")
	cmd.Dir = dirPath
	err := cmd.Run()
	if err != nil {
		return false, err.Error()
	}
	return true, ""
}
*/

func findGoFiles(dirPath string) []string {
	var goFiles []string
	files, _ := os.ReadDir(dirPath)
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".go") {
			goFiles = append(goFiles, file.Name())
		}
	}
	return goFiles
}

var packageRegex = regexp.MustCompile(`\bpackage\s+(\w+)\b`)

func GetLineColumn(content []byte, index int) (line, col int) {
	line = 1
	col = 1
	for i, ch := range content {
		if i >= index {
			break
		}
		col++
		if ch == '\n' {
			line++
			col = 1
		}
	}
	return line, col
}

func AnalyzeModule(path string, occurrences *[]*Occurrence, parser OccurrenceParser) {
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error accessing file %s: %v\n", path, err)
			return nil
		}
		if info.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}
		parser.FindOccurrences(path, occurrences)
		return nil
	})
}

func CountUniqueOccurrences(occurrences []*Occurrence) (initCount, anonymCount, execCount, pluginCount, goGenerateCount, goTestCount, unsafeCount, cgoCount, indirectCount, reflectCount, constructorCount int) {
	initOccurrences := make(map[string]struct{})
	anonymOccurrences := make(map[string]struct{})
	execOccurrences := make(map[string]struct{})
	pluginOccurrences := make(map[string]struct{})
	goGenerateOccurrences := make(map[string]struct{})
	goTestOccurrences := make(map[string]struct{})
	unsafeOccurrences := make(map[string]struct{})
	cgoOccurrences := make(map[string]struct{})
	indirectOccurrences := make(map[string]struct{})
	reflectOccurrences := make(map[string]struct{})
	constructorOccurrences := make(map[string]struct{})

	for _, occ := range occurrences {
		switch occ.AttackVector {
		case "init":
			key := fmt.Sprintf("%s:%d", occ.FilePath, occ.LineNumber)
			initOccurrences[key] = struct{}{}
		case "anonym":
			key := fmt.Sprintf("%s:%s:%d", occ.VariableName, occ.FilePath, occ.LineNumber)
			anonymOccurrences[key] = struct{}{}
		case "exec":
			key := fmt.Sprintf("%s:%s:%d", occ.MethodInvoked, occ.FilePath, occ.LineNumber)
			execOccurrences[key] = struct{}{}
		case "plugin":
			key := fmt.Sprintf("%s:%s:%d", occ.FilePath, occ.MethodInvoked, occ.LineNumber)
			pluginOccurrences[key] = struct{}{}
		case "go:generate":
			key := fmt.Sprintf("%s:%s:%d", occ.Command, occ.FilePath, occ.LineNumber)
			goGenerateOccurrences[key] = struct{}{}
		case "test":
			key := fmt.Sprintf("%s:%s:%s:%d", occ.Command, occ.FilePath, occ.MethodInvoked, occ.LineNumber)
			goTestOccurrences[key] = struct{}{}
		case "unsafe":
			key := fmt.Sprintf("%s:%s:%d", occ.MethodInvoked, occ.FilePath, occ.LineNumber) // TODO: which info to include here
			unsafeOccurrences[key] = struct{}{}
		case "cgo":
			key := fmt.Sprintf("%s:%s:%d", occ.MethodInvoked, occ.FilePath, occ.LineNumber) // TODO: which info to include here
			cgoOccurrences[key] = struct{}{}
		case "indirect":
			key := fmt.Sprintf("%s:%s:%s:%d", occ.MethodInvoked, occ.TypePassed, occ.FilePath, occ.LineNumber)
			indirectOccurrences[key] = struct{}{}
		case "reflect":
			key := fmt.Sprintf("%s:%s:%d", occ.MethodInvoked, occ.FilePath, occ.LineNumber)
			reflectOccurrences[key] = struct{}{}
		case "constructor":
			key := fmt.Sprintf("%s:%s:%d", occ.FilePath, occ.Pattern, occ.LineNumber)
			constructorOccurrences[key] = struct{}{}
		}
	}

	return len(initOccurrences), len(anonymOccurrences), len(execOccurrences), len(pluginOccurrences), len(goGenerateOccurrences), len(goTestOccurrences), len(unsafeOccurrences), len(cgoOccurrences), len(indirectOccurrences), len(reflectOccurrences), len(constructorOccurrences)
}

type OccurrenceJSON struct {
	Type          string `json:"Type,omitempty"`
	FilePath      string `json:"FilePath,omitempty"`
	LineNumber    int    `json:"LineNumber,omitempty"`
	MethodInvoked string `json:"MethodInvoked,omitempty"`
	TypePassed    string `json:"TypePassed,omitempty"`
	VariableName  string `json:"VariableName,omitempty"`
	Command       string `json:"Command,omitempty"`
	Pattern       string `json:"Pattern,omitempty"`
}

func PrintOccurrences(occurrences []*Occurrence) {
	var result []OccurrenceJSON
	for _, occ := range occurrences {
		occJSON := OccurrenceJSON{
			Type:          occ.AttackVector,
			FilePath:      occ.FilePath,
			LineNumber:    occ.LineNumber,
			MethodInvoked: occ.MethodInvoked,
			TypePassed:    occ.TypePassed,
			VariableName:  occ.VariableName,
			Command:       occ.Command,
			Pattern:       occ.Pattern,
		}
		result = append(result, occJSON)
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	fmt.Println("Vector Occurrences:")
	fmt.Println(string(jsonData))
}

func PrintDependencies(dependencies []Dependency) {
	jsonDependencies, err := json.MarshalIndent(dependencies, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}
	fmt.Println(string(jsonDependencies))
}

// Function to render a progress bar on the console
func updateProgressBar(current, total int) {
	width := 50 // Width of the progress bar
	progress := float64(current) / float64(total)
	hashes := int(progress * float64(width))
	fmt.Printf("\r[%-*s] %.2f%%", width, strings.Repeat("#", hashes), progress*100)
}

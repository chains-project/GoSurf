package analysis

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
)

type Occurrence struct {
	PackageName   string
	AttackVector  string
	FilePath      string
	LineNumber    int
	VariableName  string // for anonymous functions
	Command       string // for go:generate directive
	MethodInvoked string // for indirect, exec, plugin, cgo
	TypePassed    string // for indirect
	Pattern       string // for constructors
}

type OccurrenceJSON struct {
	PackageName   string `json:"PackageName,omitempty"`
	Type          string `json:"Type,omitempty"`
	FilePath      string `json:"FilePath,omitempty"`
	LineNumber    int    `json:"LineNumber,omitempty"`
	MethodInvoked string `json:"MethodInvoked,omitempty"`
	TypePassed    string `json:"TypePassed,omitempty"`
	VariableName  string `json:"VariableName,omitempty"`
	Command       string `json:"Command,omitempty"`
	Pattern       string `json:"Pattern,omitempty"`
}

type Dependency struct {
	Name string
	Path string
}

type OccurrenceParser interface {
	FindOccurrences(path string, packageName string, occurrences *[]*Occurrence)
}

var pkgAsmFunctions []string

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
	goFiles := findFiles(".go", dirPath)
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

func findFiles(suffix string, dirPath string) []string {
	var files []string
	dirFiles, _ := os.ReadDir(dirPath)
	for _, file := range dirFiles {
		if !file.IsDir() && strings.HasSuffix(file.Name(), suffix) {
			files = append(files, file.Name())
		}
	}
	return files
}

func pkgContainsAsm(dirPath string) (bool, []string) {
	var asmSuffixes = []string{".s", ".S", ".sx"}
	var files []string

	for _, suffix := range asmSuffixes {
		files = append(files, findFiles(suffix, dirPath)...)
	}
	if len(files) == 0 {
		// No assembly files in package
		return false, nil
	}

	// Find all assembly function signatures in pkg ('TEXT ·' pattern) 
	var signatureRegex = regexp.MustCompile(`TEXT\s+·[A-Za-z1-9]+\w*`)
	var signatures []string

	for _, file := range files {
		filePath := filepath.Join(dirPath, file)
		content, err := os.ReadFile(filePath)
		if err != nil {
			continue
		}
		match := signatureRegex.FindString(string(content))
		if match != "" {
			funSig := strings.Split(match, "·")[1]
			signatures = append(signatures, funSig)
		}
	}
	return true, signatures
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

func AnalyzePackage(dep Dependency, occurrences *[]*Occurrence, parser OccurrenceParser) {
	files, err := os.ReadDir(dep.Path)
	if err != nil {
		fmt.Printf("Error accessing directory %s: %v\n", dep.Path, err)
		return
	}

	// If assembly parser, get set the assembly function definitions in the package
	currentParser := reflect.TypeOf(parser)
	asmParser := reflect.TypeOf(AssemblyParser{})
	
	if currentParser == asmParser {
		ok, funSigs := pkgContainsAsm(dep.Path)
		pkgAsmFunctions = funSigs
		if !ok {
			// Avoid running assembly parser in package without assembly
			return
		}
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".go") {
			continue
		}
		path := filepath.Join(dep.Path, file.Name())
		parser.FindOccurrences(path, dep.Name, occurrences)
	}
}

func CountUniqueOccurrences(occurrences []*Occurrence) (initCount, anonymCount, execCount, pluginCount, goGenerateCount, goTestCount, unsafeCount, cgoCount, indirectCount, reflectCount, constructorCount, assemblyCount int) {
	initOccurrences := make(map[string]struct{})
	globalVarOccurrences := make(map[string]struct{})
	execOccurrences := make(map[string]struct{})
	pluginOccurrences := make(map[string]struct{})
	goGenerateOccurrences := make(map[string]struct{})
	goTestOccurrences := make(map[string]struct{})
	unsafeOccurrences := make(map[string]struct{})
	cgoOccurrences := make(map[string]struct{})
	indirectOccurrences := make(map[string]struct{})
	reflectOccurrences := make(map[string]struct{})
	constructorOccurrences := make(map[string]struct{})
	assemblyOccurrences := make(map[string]struct{})

	for _, occ := range occurrences {
		switch occ.AttackVector {
		case "init":
			key := fmt.Sprintf("%s:%d", occ.FilePath, occ.LineNumber)
			initOccurrences[key] = struct{}{}
		case "global":
			key := fmt.Sprintf("%s:%s:%d", occ.VariableName, occ.FilePath, occ.LineNumber)
			globalVarOccurrences[key] = struct{}{}
		case "exec":
			key := fmt.Sprintf("%s:%s:%d", occ.MethodInvoked, occ.FilePath, occ.LineNumber)
			execOccurrences[key] = struct{}{}
		case "plugin":
			key := fmt.Sprintf("%s:%s:%d", occ.FilePath, occ.MethodInvoked, occ.LineNumber)
			pluginOccurrences[key] = struct{}{}
		case "generate":
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
		case "assembly":
			key := fmt.Sprintf("%s:%s:%d", occ.MethodInvoked, occ.FilePath, occ.LineNumber)
			assemblyOccurrences[key] = struct{}{}
		}
	}

	return len(initOccurrences), len(globalVarOccurrences), len(execOccurrences), len(pluginOccurrences), len(goGenerateOccurrences), len(goTestOccurrences), len(unsafeOccurrences), len(cgoOccurrences), len(indirectOccurrences), len(reflectOccurrences), len(constructorOccurrences), len(assemblyOccurrences)
}

func PrintOccurrences(occurrences []*Occurrence) {
	var result []OccurrenceJSON
	for _, occ := range occurrences {
		occJSON := OccurrenceJSON{
			PackageName:   occ.PackageName,
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
	/*
		jsonDependencies, err := json.MarshalIndent(dependencies, "", "  ")
		if err != nil {
			fmt.Println("Error marshaling JSON:", err)
			return
		}
		fmt.Println(string(jsonDependencies))
	*/
	for _, dep := range dependencies {
		fmt.Println(dep.Path)
	}
}

// Function to render a progress bar on the console
func updateProgressBar(current, total int) {
	width := 50 // Width of the progress bar
	progress := float64(current) / float64(total)
	hashes := int(progress * float64(width))
	fmt.Printf("\r[%-*s] %.2f%%", width, strings.Repeat("#", hashes), progress*100)
}

package analysis

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"regexp"
	"strings"
)

type InitFuncParser struct{}
type AnonymFuncParser struct{}
type OsExecParser struct{}
type PluginParser struct{}

// Parser for Anonym Function analysis
func (p InitFuncParser) FindOccurrences(path string, occurrences *[]*Occurrence) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		fmt.Printf("Error parsing file %s: %v\n", path, err)
		return
	}

	for _, decl := range node.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Name.Name != "init" {
			continue
		}

		*occurrences = append(*occurrences, &Occurrence{
			Type:       "init",
			FilePath:   path,
			LineNumber: fset.Position(fn.Pos()).Line,
		})
	}
}

// Parser for init() Function analysis
func (p AnonymFuncParser) FindOccurrences(path string, occurrences *[]*Occurrence) {
	fileContents, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", path, err)
		return
	}

	pattern := `var\s+(\w+)\s*(\w*)\s*=\s*func\(\)\s*(\w*)\s*{[^}]*}\(\)`
	re := regexp.MustCompile(pattern)

	matches := re.FindAllStringSubmatchIndex(string(fileContents), -1)
	if len(matches) > 0 {
		for _, match := range matches {
			startLine, _ := GetLineColumn(fileContents, match[0])
			variableName := strings.TrimSpace(string(fileContents[match[2]:match[3]]))
			*occurrences = append(*occurrences, &Occurrence{
				Type:         "anonym",
				VariableName: variableName,
				FilePath:     path,
				LineNumber:   startLine,
			})
		}
	}
}

// Parser for os/exec Function analysis
func (p OsExecParser) FindOccurrences(path string, occurrences *[]*Occurrence) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		fmt.Printf("Error parsing file %s: %v\n", path, err)
		return
	}

	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.CallExpr:
			fun, ok := x.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}

			pkg, ok := fun.X.(*ast.Ident)
			if !ok || pkg.Name != "exec" {
				return true
			}

			// Add this if you want to check for specific functions
			// if strings.HasPrefix(fun.Sel.Name, "Command") {

			// Check for all functions within the os/exec package
			*occurrences = append(*occurrences, &Occurrence{
				Type:       "exec",
				Function:   fun.Sel.Name,
				FilePath:   path,
				LineNumber: fset.Position(x.Pos()).Line,
			})
		}
		return true
	})
}

// Parser for Go plugin usage
func (p PluginParser) FindOccurrences(path string, occurrences *[]*Occurrence) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		fmt.Printf("Error parsing file %s: %v\n", path, err)
		return
	}

	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.CallExpr:
			fun, ok := x.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}

			pkg, ok := fun.X.(*ast.Ident)
			if !ok || pkg.Name != "plugin" {
				return true
			}

			// Check for plugin.Open function
			if fun.Sel.Name == "Open" {
				*occurrences = append(*occurrences, &Occurrence{
					Type:       "plugin",
					Function:   fun.Sel.Name,
					FilePath:   path,
					LineNumber: fset.Position(x.Pos()).Line,
				})
			}
		}
		return true
	})
}

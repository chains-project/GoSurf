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
type ExecParser struct{}
type PluginParser struct{}
type GoGenerateParser struct{}
type UnsafeParser struct{}
type CgoParser struct{}

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

type execFuncInfo struct {
	pkgName   string
	funcNames []string
}

// Add here exec functions to check for exec analysis
var execFuncs = []execFuncInfo{
	{"syscall", []string{"Exec", "ForkExec", "StartProcess"}},
	{"exec", []string{"Command", "CommandContext"}},
	{"os", []string{"StartProcess"}},
}

// Parser for exec function analysis
func (p ExecParser) FindOccurrences(path string, occurrences *[]*Occurrence) {
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
			if !ok {
				return true
			}

			for _, execFunc := range execFuncs {
				if pkg.Name == execFunc.pkgName {
					for _, funcName := range execFunc.funcNames {
						if fun.Sel.Name == funcName {
							*occurrences = append(*occurrences, &Occurrence{
								Type:       "exec",
								Function:   pkg.Name + "." + fun.Sel.Name,
								FilePath:   path,
								LineNumber: fset.Position(x.Pos()).Line,
							})
							break
						}
					}
				}
			}

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

// Parser for go:generate directive analysis
func (p GoGenerateParser) FindOccurrences(path string, occurrences *[]*Occurrence) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		fmt.Printf("Error parsing file %s: %v\n", path, err)
		return
	}

	for _, cg := range node.Comments {
		for _, c := range cg.List {
			if strings.HasPrefix(c.Text, "//go:generate") {
				*occurrences = append(*occurrences, &Occurrence{
					Type:       "go:generate",
					Command:    strings.TrimPrefix(c.Text, "//go:generate "),
					FilePath:   path,
					LineNumber: fset.Position(c.Pos()).Line,
				})
			}
		}
	}
}

// Parser for unsafe pointer usage
func (p UnsafeParser) FindOccurrences(path string, occurrences *[]*Occurrence) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		fmt.Printf("Error parsing file %s: %v\n", path, err)
		return
	}
	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.CallExpr:
			if sel, ok := x.Fun.(*ast.SelectorExpr); ok {
				if pkg, ok := sel.X.(*ast.Ident); ok && pkg.Name == "unsafe" && sel.Sel.Name == "Pointer" {
					*occurrences = append(*occurrences, &Occurrence{
						Type:       "unsafe",
						FilePath:   path,
						LineNumber: fset.Position(x.Pos()).Line,
					})
				}
			}
		}
		return true
	})
}

func (p CgoParser) FindOccurrences(path string, occurrences *[]*Occurrence) {

	// todo
	fmt.Println("do something")
}

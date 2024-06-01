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
type IndirectParser struct{}
type ReflectParser struct{}

// Parser for init() Function analysis
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
			AttackVector: "init",
			FilePath:     path,
			LineNumber:   fset.Position(fn.Pos()).Line,
		})
	}
}

// Parser for anonymous functions analysis
func (p AnonymFuncParser) FindOccurrences(path string, occurrences *[]*Occurrence) {
	fileContents, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", path, err)
		return
	}

	//pattern := `var\s+(\w+)(?:\s+\w+)?\s*=\s*func\(\)(?:\s*\w+)?\s*{\s*[^}]*\s*}\(\)`
	pattern := `var\s+(\w+)\s*(\w*)\s*=\s*func\(\)\s*(\w*)\s*{[^}]*}\(\)`
	re := regexp.MustCompile(pattern)

	matches := re.FindAllStringSubmatchIndex(string(fileContents), -1)
	if len(matches) > 0 {
		for _, match := range matches {
			startLine, _ := GetLineColumn(fileContents, match[0])
			variableName := strings.TrimSpace(string(fileContents[match[2]:match[3]]))
			*occurrences = append(*occurrences, &Occurrence{
				AttackVector: "anonym",
				FilePath:     path,
				LineNumber:   startLine,
				VariableName: variableName,
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
								AttackVector:  "exec",
								FilePath:      path,
								LineNumber:    fset.Position(x.Pos()).Line,
								MethodInvoked: pkg.Name + "." + fun.Sel.Name,
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
					AttackVector:  "plugin",
					FilePath:      path,
					LineNumber:    fset.Position(x.Pos()).Line,
					MethodInvoked: fun.Sel.Name,
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
					AttackVector: "go:generate",
					FilePath:     path,
					LineNumber:   fset.Position(c.Pos()).Line,
					Command:      strings.TrimPrefix(c.Text, "//go:generate "),
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
						AttackVector:  "unsafe",
						FilePath:      path,
						LineNumber:    fset.Position(x.Pos()).Line,
						MethodInvoked: "unsafe.Pointer",
					})
				}
			}
		}
		return true
	})
}

// Parser for Cgo usage
func (p CgoParser) FindOccurrences(path string, occurrences *[]*Occurrence) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		fmt.Printf("Error parsing file %s: %v\n", path, err)
		return
	}

	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.CallExpr:
			sel, ok := x.Fun.(*ast.SelectorExpr)

			if !ok {
				return true
			}

			if pkg, ok := sel.X.(*ast.Ident); ok && pkg.Name == "C" {
				*occurrences = append(*occurrences, &Occurrence{
					AttackVector:  "cgo",
					FilePath:      path,
					LineNumber:    fset.Position(x.Pos()).Line,
					MethodInvoked: "C." + sel.Sel.Name,
				})
			}
		}
		return true
	})
}

// Parser for indirect method calls through interface
func (p IndirectParser) FindOccurrences(path string, occurrences *[]*Occurrence) {

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		fmt.Printf("Error parsing file %s: %v\n", path, err)
		return
	}

	methods := make(map[string][]string)

	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			if x.Recv != nil {
				receiverType := fmt.Sprint(x.Recv.List[0].Type)
				methods[x.Name.Name] = append(methods[x.Name.Name], receiverType)
			}
		}
		return true
	})

	polymorphicMethods := make(map[string]struct{})
	for name, receiverTypes := range methods {
		receiverTypeSet := make(map[string]struct{})
		for _, t := range receiverTypes {
			receiverTypeSet[t] = struct{}{}
		}
		if len(receiverTypeSet) > 1 {
			polymorphicMethods[name] = struct{}{}
		}
	}

	/* To save polymorphic methods definitions
	for method, receiverTypes := range polymorphicMethods {
		for _, receiverType := range receiverTypes {
			*occurrences = append(*occurrences, &Occurrence{
				Type:            	"interface",
				VariablesPassed:   	receiverType,
				MethodInvoked: 		method,
				FilePath:        	path,
			})
		}
	}
	*/

	// Find invocations of polymorphic methods
	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.CallExpr:
			fun, ok := x.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}

			_, isPolymorphic := polymorphicMethods[fun.Sel.Name]
			if isPolymorphic {
				receiverType := ""
				if expr, ok := x.Fun.(*ast.SelectorExpr); ok {
					if ident, ok := expr.X.(*ast.Ident); ok {
						receiverType = ident.Name
					}
				}
				*occurrences = append(*occurrences, &Occurrence{
					AttackVector:  "indirect",
					FilePath:      path,
					LineNumber:    fset.Position(x.Pos()).Line,
					MethodInvoked: fun.Sel.Name,
					TypePassed:    receiverType,
				})
			}
		}
		return true
	})
}

func (p ReflectParser) FindOccurrences(path string, occurrences *[]*Occurrence) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		fmt.Printf("Error parsing file %s: %v\n", path, err)
		return
	}

	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.ImportSpec:
			if pkg := x.Path.Value; pkg == `"reflect"` {
				*occurrences = append(*occurrences, &Occurrence{
					AttackVector: "reflect",
					FilePath:     path,
					LineNumber:   fset.Position(x.Pos()).Line})
				return false
			}
		}
		return true
	})
}

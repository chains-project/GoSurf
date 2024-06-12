package analysis

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type InitFuncParser struct{}
type GlobalVarParser struct{}
type ExecParser struct{}
type PluginParser struct{}
type GoGenerateParser struct{}
type GoTestParser struct{}
type UnsafeParser struct{}
type CgoParser struct{}
type InterfaceParser struct{}
type ReflectParser struct{}
type ConstructorParser struct{}
type AssemblyParser struct{}

// Parser for init() Function analysis
func (p InitFuncParser) FindOccurrences(path string, packageName string, occurrences *[]*Occurrence) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, nil, parser.AllErrors)
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
			PackageName:  packageName,
			AttackVector: "init",
			FilePath:     path,
			LineNumber:   fset.Position(fn.Pos()).Line,
		})
	}
}

// Parser for global var initialization with functions
func (p GlobalVarParser) FindOccurrences(path string, packageName string, occurrences *[]*Occurrence) {

	fileContents, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", path, err)
		return
	}

	// Global variable initialization with normal functions
	normal_pattern := `var\s+(\w+)\s*(\w*)\s*=\s*(\w+)\(\)`
	normal_re := regexp.MustCompile(normal_pattern)
	normal_matches := normal_re.FindAllStringSubmatchIndex(string(fileContents), -1)
	if len(normal_matches) > 0 {
		for _, match := range normal_matches {
			startLine, _ := GetLineColumn(fileContents, match[0])
			variableName := strings.TrimSpace(string(fileContents[match[2]:match[3]]))
			funcName := strings.TrimSpace(string(fileContents[match[6]:match[7]]))

			// Skip if the function name is "func"
			if funcName == "func" {
				continue
			}

			*occurrences = append(*occurrences, &Occurrence{
				PackageName:   packageName,
				AttackVector:  "global",
				FilePath:      path,
				LineNumber:    startLine,
				VariableName:  variableName,
				MethodInvoked: funcName + "()",
			})
		}
	}

	// Global variable initialization with anonymous functions
	anonym_pattern := `var\s+(\w+)\s*(\w*)\s*=\s*func\(\)\s*(\w*)\s*{[^}]*}\(\)`
	anonym_re := regexp.MustCompile(anonym_pattern)
	anonym_matches := anonym_re.FindAllStringSubmatchIndex(string(fileContents), -1)
	if len(anonym_matches) > 0 {
		for _, match := range anonym_matches {
			startLine, _ := GetLineColumn(fileContents, match[0])
			variableName := strings.TrimSpace(string(fileContents[match[2]:match[3]]))
			*occurrences = append(*occurrences, &Occurrence{
				AttackVector:  "global",
				FilePath:      path,
				LineNumber:    startLine,
				VariableName:  variableName,
				MethodInvoked: "anonym func",
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
func (p ExecParser) FindOccurrences(path string, packageName string, occurrences *[]*Occurrence) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, nil, parser.AllErrors)
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
								PackageName:   packageName,
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
func (p PluginParser) FindOccurrences(path string, packageName string, occurrences *[]*Occurrence) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, nil, parser.AllErrors)
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
					PackageName:   packageName,
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
func (p GoGenerateParser) FindOccurrences(path string, packageName string, occurrences *[]*Occurrence) {
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
					PackageName:  packageName,
					AttackVector: "generate",
					FilePath:     path,
					LineNumber:   fset.Position(c.Pos()).Line,
					Command:      strings.TrimPrefix(c.Text, "//go:generate "),
				})
			}
		}
	}
}

// Parser for Test functions (prefix: Test, Benchmark, Example) analysis
func (p GoTestParser) FindOccurrences(path string, packageName string, occurrences *[]*Occurrence) {

	if strings.HasSuffix(path, "_test.go") {
		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, path, nil, parser.AllErrors)
		if err != nil {
			fmt.Printf("Error parsing file %s: %v\n", path, err)
			return
		}

		ast.Inspect(node, func(n ast.Node) bool {
			fn, ok := n.(*ast.FuncDecl)
			if !ok {
				return true
			}

			funcName := fn.Name.Name
			if strings.HasPrefix(funcName, "Test") || strings.HasPrefix(funcName, "Benchmark") || strings.HasPrefix(funcName, "Example") || strings.HasPrefix(funcName, "Fuzz") {

				filePath := filepath.Join(path, node.Name.Name)
				*occurrences = append(*occurrences, &Occurrence{
					PackageName:   packageName,
					AttackVector:  "test",
					FilePath:      filePath,
					LineNumber:    fset.Position(fn.Pos()).Line,
					MethodInvoked: funcName,
				})
			}

			return true
		})
	}

}

// Parser for unsafe pointer usage
func (p UnsafeParser) FindOccurrences(path string, packageName string, occurrences *[]*Occurrence) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, nil, parser.AllErrors)
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
						PackageName:   packageName,
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
func (p CgoParser) FindOccurrences(path string, packageName string, occurrences *[]*Occurrence) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, nil, parser.AllErrors)
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
					PackageName:   packageName,
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

// Parser for indirect method invocations throguh Interfaces
func (p InterfaceParser) FindOccurrences(path string, packageName string, occurrences *[]*Occurrence) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, nil, parser.AllErrors)
	if err != nil {
		fmt.Printf("Error parsing file %s: %v\n", path, err)
		return
	}

	methods := make(map[string][]string)
	interfaceMethods := make(map[string]struct{})

	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			if x.Recv != nil {
				receiverType := fmt.Sprint(x.Recv.List[0].Type)
				methods[x.Name.Name] = append(methods[x.Name.Name], receiverType)
			} else if x.Type.Results != nil && len(x.Type.Results.List) > 0 {
				// Check if the function is defined on an interface type
				if _, ok := x.Type.Results.List[0].Type.(*ast.InterfaceType); ok {
					interfaceMethods[x.Name.Name] = struct{}{}
				}
			}
		}
		return true
	})

	// Mark all interface methods as polymorphic
	polymorphicMethods := make(map[string]struct{})
	for name, receiverTypes := range methods {
		receiverTypeSet := make(map[string]struct{})
		for _, t := range receiverTypes {
			receiverTypeSet[t] = struct{}{}
		}
		polymorphicMethods[name] = struct{}{}
	}

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
					PackageName:   packageName,
					AttackVector:  "interface",
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

func (p ReflectParser) FindOccurrences(path string, packageName string, occurrences *[]*Occurrence) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, nil, parser.AllErrors)
	if err != nil {
		fmt.Printf("Error parsing file %s: %v\n", path, err)
		return
	}

	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.ImportSpec:
			if pkg := x.Path.Value; pkg == `"reflect"` {
				*occurrences = append(*occurrences, &Occurrence{
					PackageName:  packageName,
					AttackVector: "reflect",
					FilePath:     path,
					LineNumber:   fset.Position(x.Pos()).Line})
				return false
			}
		}
		return true
	})
}

func (p ConstructorParser) FindOccurrences(path string, packageName string, occurrences *[]*Occurrence) {

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, nil, parser.AllErrors)
	if err != nil {
		fmt.Printf("Error parsing file %s: %v\n", path, err)
		return
	}

	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.CallExpr:

			switch fun := x.Fun.(type) {

			// Check for factory function invocations
			case *ast.SelectorExpr:
				if strings.HasPrefix(fun.Sel.Name, "New") {
					*occurrences = append(*occurrences, &Occurrence{
						PackageName:   packageName,
						AttackVector:  "constructor",
						FilePath:      path,
						LineNumber:    fset.Position(x.Pos()).Line,
						MethodInvoked: fun.Sel.Name,
						Pattern:       "factory function",
					})
				}

			// Check for `New` function invocations
			case *ast.Ident:
				if fun.Name == "New" {
					*occurrences = append(*occurrences, &Occurrence{
						PackageName:   packageName,
						AttackVector:  "constructor",
						FilePath:      path,
						LineNumber:    fset.Position(x.Pos()).Line,
						MethodInvoked: "New",
						Pattern:       "New() function",
					})
				}
			}

		}
		return true
	})
}

// Parser for Assembly function use
func (p AssemblyParser) FindOccurrences(path string, packageName string, occurrences *[]*Occurrence) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, nil, parser.AllErrors)
	if err != nil {
		fmt.Printf("Error parsing file %s: %v\n", path, err)
		return
	}

	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.CallExpr:
			if fun, ok := x.Fun.(*ast.Ident); ok {
				for _, funSig := range pkgAsmFunctions {
					if fun.Name == funSig {
						*occurrences = append(*occurrences, &Occurrence{
							PackageName:   packageName,
							AttackVector:  "assembly",
							FilePath:      path,
							LineNumber:    fset.Position(x.Pos()).Line,
							MethodInvoked: fun.Name,
						})
						break
					}
				}
			}
		}
		return true
	})
}

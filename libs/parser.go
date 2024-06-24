package analysis

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
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

// Parser for init() function declarations.
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

// Parser for global variable declarations.
func (p GlobalVarParser) FindOccurrences(path string, packageName string, occurrences *[]*Occurrence) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, nil, parser.AllErrors)
	if err != nil {
		fmt.Printf("Error parsing file %s: %v\n", path, err)
		return
	}

	// common pattern: *ast.GenDecl.Specs.*ast.ValueSpec.Values[0].CallExpr
	for _, decl := range node.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok || len(gd.Specs) == 0 {
			continue
		}
		//t
		for _, spec := range gd.Specs {
			val, ok := spec.(*ast.ValueSpec)

			//val, ok := gd.Specs[0].(*ast.ValueSpec) // constant or variable declaration
			if !ok || len(val.Values) == 0 {
				continue
			}

			if exp, ok := val.Values[0].(*ast.CallExpr); ok { // function call
				name := val.Names[0].Name
				switch x := exp.Fun.(type) {

				case *ast.Ident:
					*occurrences = append(*occurrences, &Occurrence{
						PackageName:   packageName,
						AttackVector:  "global",
						FilePath:      path,
						LineNumber:    fset.Position(x.Pos()).Line,
						VariableName:  name,
						MethodInvoked: x.Name + "()",
					})
					continue

				case *ast.SelectorExpr:
					*occurrences = append(*occurrences, &Occurrence{
						PackageName:   packageName,
						AttackVector:  "global",
						FilePath:      path,
						LineNumber:    fset.Position(x.Pos()).Line,
						VariableName:  name,
						MethodInvoked: x.Sel.Name + "()",
					})
					continue

				case *ast.FuncLit:
					*occurrences = append(*occurrences, &Occurrence{
						AttackVector:  "global",
						FilePath:      path,
						LineNumber:    fset.Position(x.Pos()).Line,
						VariableName:  name,
						MethodInvoked: "anonym func",
					})
					continue

				}

			}
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
		if x, ok := n.(*ast.CallExpr); ok {
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
		if x, ok := n.(*ast.CallExpr); ok {
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

// Parser for go:generate directive analysis.
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

// Parser for Test functions (prefix: Test, Benchmark, Example) analysis.
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

// Parser for unsafe pointer usage.
func (p UnsafeParser) FindOccurrences(path string, packageName string, occurrences *[]*Occurrence) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, nil, parser.AllErrors)
	if err != nil {
		fmt.Printf("Error parsing file %s: %v\n", path, err)
		return
	}

	ast.Inspect(node, func(n ast.Node) bool {
		if x, ok := n.(*ast.CallExpr); ok {
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

// Parser for Cgo usage.
func (p CgoParser) FindOccurrences(path string, packageName string, occurrences *[]*Occurrence) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, nil, parser.AllErrors)
	if err != nil {
		fmt.Printf("Error parsing file %s: %v\n", path, err)
		return
	}

	ast.Inspect(node, func(n ast.Node) bool {
		if x, ok := n.(*ast.CallExpr); ok {
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

// Parser for indirect method invocations throguh Interfaces.
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
		if x, ok := n.(*ast.FuncDecl); ok {
			if x.Recv != nil && len(x.Recv.List) != 0 {
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
		if x, ok := n.(*ast.CallExpr); ok {
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

// Parser for imports of reflect package.
func (p ReflectParser) FindOccurrences(path string, packageName string, occurrences *[]*Occurrence) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, nil, parser.AllErrors)
	if err != nil {
		fmt.Printf("Error parsing file %s: %v\n", path, err)
		return
	}

	for _, decl := range node.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok || len(gd.Specs) == 0 {
			continue
		}
		for _, spec := range gd.Specs {
			im, ok := spec.(*ast.ImportSpec)
			if !ok {
				continue
			}
			if pkg := im.Path.Value; pkg == `"reflect"` {
				*occurrences = append(*occurrences, &Occurrence{
					PackageName:  packageName,
					AttackVector: "reflect",
					FilePath:     path,
					LineNumber:   fset.Position(im.Pos()).Line})
				break
			}
		}
	}
}

// Parser for constructor usage.
func (p ConstructorParser) FindOccurrences(path string, packageName string, occurrences *[]*Occurrence) {

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, nil, parser.AllErrors)
	if err != nil {
		fmt.Printf("Error parsing file %s: %v\n", path, err)
		return
	}

	ast.Inspect(node, func(n ast.Node) bool {
		if x, ok := n.(*ast.CallExpr); ok {
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

// Parser for Assembly function usage.
func (p AssemblyParser) FindOccurrences(path string, packageName string, occurrences *[]*Occurrence) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, nil, parser.AllErrors)
	if err != nil {
		fmt.Printf("Error parsing file %s: %v\n", path, err)
		return
	}

	ast.Inspect(node, func(n ast.Node) bool {
		if x, ok := n.(*ast.CallExpr); ok {
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

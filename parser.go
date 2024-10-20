package gogo

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

const GOGOIMPORTPATH = "github.com/morganhein/gogo"

func parse(src string) (Context, error) {
	// Parse the source code
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		panic(err)
	}

	// Find the import alias for the gogo package
	gogoAlias, err := getGoGoImportName(file)
	if err != nil {
		panic(err)
	}

	// Traverse the AST and find the desired function
	ast.Inspect(file, func(node ast.Node) bool {
		funcDecl, ok := node.(*ast.FuncDecl)
		if !ok || funcDecl.Name.Name != "UpperCaseName" {
			return true
		}

		// Extract information from ctx.<method> calls
		for _, stmt := range funcDecl.Body.List {
			// determine if the root exprStmt is "c
			exprStmt, ok := stmt.(*ast.ExprStmt)
			if !ok {
				continue
			}
			isGoGoCtx := determineIfRootIsGoGoCtx(gogoAlias, exprStmt)
			if !isGoGoCtx {
				continue
			}
		}

		return false
	})
	return nil, nil
}

// parseGoGoCtx extracts GoGoContext information from the AST
func parseGoGoCtx(exprStmt *ast.ExprStmt) {
	//switch selExpr.Sel.Name {
	//case "SetDescription":
	//	arg := callExpr.Args[0].(*ast.BasicLit).Value
	//	fmt.Printf("Description: %s\n", arg)
	//case "SetHelp":
	//	arg := callExpr.Args[0].(*ast.BasicLit).Value
	//	fmt.Printf("Help: %s\n", arg)
	//case "Argument":
	//	parseArgument(exprStmt)
	//default:
	//	fmt.Printf("Could not determine: %v", selExpr.Sel.Name)
	//}
	fmt.Println(exprStmt.Pos())
}

// parseArgument extracts GoGoArgument information from the AST
func parseArgument(exprStmt *ast.ExprStmt) {
	callExpr := exprStmt.X.(*ast.CallExpr)
	optionExpr := callExpr.Args[0]
	optionName := ""

	switch expr := optionExpr.(type) {
	case *ast.Ident:
		optionName = expr.Name
	case *ast.SelectorExpr:
		optionName = expr.Sel.Name
	}

	fmt.Printf("Argument: %s\n", optionName)

	// Check for chained method calls
	for _, stmt := range callExpr.Args {
		callExpr, ok := stmt.(*ast.CallExpr)
		if !ok {
			continue
		}

		selExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
		if !ok {
			continue
		}

		switch selExpr.Sel.Name {
		case "Description":
			arg := callExpr.Args[0].(*ast.BasicLit).Value
			fmt.Printf("  Description: %s\n", arg)
		case "AllowedValues":
			var allowedValues []string
			for _, arg := range callExpr.Args {
				allowedValues = append(allowedValues, arg.(*ast.BasicLit).Value)
			}
			fmt.Printf("  AllowedValues: %v\n", allowedValues)
		case "RestrictedValues":
			var restrictedValues []string
			for _, arg := range callExpr.Args {
				restrictedValues = append(restrictedValues, arg.(*ast.BasicLit).Value)
			}
			fmt.Printf("  RestrictedValues: %v\n", restrictedValues)
		case "Default":
			arg := callExpr.Args[0].(*ast.BasicLit).Value
			fmt.Printf("  Default: %s\n", arg)
		}
	}
}

// getGoGoFunctions parses the AST and finds all the functions
// that are valid GoGoFunctions/entrypoints
func getGoGoFunctions(gogoAlias string, file *ast.File) []*ast.FuncDecl {
	var validFunctions []*ast.FuncDecl

	for _, decl := range file.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		// Check if the function is exported
		if !funcDecl.Name.IsExported() {
			continue
		}

		// Check if the function has at least one parameter
		if len(funcDecl.Type.Params.List) == 0 {
			continue
		}

		// Check if the first parameter is of type gogo.Context
		firstParam := funcDecl.Type.Params.List[0]
		if len(firstParam.Names) != 1 {
			continue
		}

		firstParamType, ok := firstParam.Type.(*ast.SelectorExpr)
		if !ok {
			continue
		}

		if firstParamType.X.(*ast.Ident).Name != gogoAlias || firstParamType.Sel.Name != "Context" {
			continue
		}

		// Check if the return value is either empty or an error
		if len(funcDecl.Type.Results.List) > 1 {
			continue
		}

		if len(funcDecl.Type.Results.List) == 1 {
			returnType, ok := funcDecl.Type.Results.List[0].Type.(*ast.Ident)
			if !ok || returnType.Name != "error" {
				continue
			}
		}

		validFunctions = append(validFunctions, funcDecl)
	}

	return validFunctions
}

// getGoGoImportName finds the alias or default name of the GoGo import
func getGoGoImportName(file *ast.File) (string, error) {
	for _, imp := range file.Imports {
		if imp.Path.Value == fmt.Sprintf("\"%s\"", GOGOIMPORTPATH) {
			if imp.Name != nil {
				return imp.Name.Name, nil
			}
			return "gogo", nil
		}
	}
	return "", errors.New("gogo import not found")
}

// determineIfRootIsGoGoCtx determines if the given exprStmnt, or any of it's children,
// are a valid GoGoContext method call. We have to walk down the AST to determine if
// the root is a GoGoContext method call because the root may be a chained method call.
func determineIfRootIsGoGoCtx(gogoAlias string, exprStmt *ast.ExprStmt) bool {
	if exprStmt.X == nil {
		return false
	}
	callExpr, ok := exprStmt.X.(*ast.CallExpr)
	if !ok {
		return false
	}
	selExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	fmt.Println(selExpr.Pos())

	return true
}

// getFuncByName returns a function declaration by name
func getFuncByName(funcDecls []*ast.FuncDecl, name string) *ast.FuncDecl {
	for _, funcDecl := range funcDecls {
		if funcDecl.Name.Name == name {
			return funcDecl
		}
	}
	return nil
}

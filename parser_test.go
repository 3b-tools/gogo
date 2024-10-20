package gogo

import (
	"go/ast"
	"go/parser"
	"go/token"
	"slices"
	"testing"
)

// test function that asserts the functions found match the expected names
func validateFunctions(names []string, funcDecls []*ast.FuncDecl) bool {
	// Create a slice to store the names of funcDecls
	funcDeclNames := make([]string, len(funcDecls))
	for i, funcDecl := range funcDecls {
		funcDeclNames[i] = funcDecl.Name.Name
	}

	slices.Sort(funcDeclNames)
	slices.Sort(names)
	diffCount := slices.Compare(funcDeclNames, names)
	return diffCount == 0
}

func TestParseImports(t *testing.T) {
	tests := []struct {
		name           string
		filePath       string
		expectedImport string
		err            bool
	}{
		{
			name:     "no imports",
			filePath: "testdata/broken.go",
			err:      true,
		},
		{
			name:           "basic",
			filePath:       "testdata/basic.go",
			expectedImport: "gogo",
		},
		{
			name:           "aliased",
			filePath:       "testdata/aliased.go",
			expectedImport: "gogo2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse the source code
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, tt.filePath, nil, 0)
			if err != nil {
				t.Fatalf("error parsing file: %v", err)
			}

			importName, err := getGoGoImportName(file)
			if (err != nil) != tt.err {
				t.Fatal("expected an error")
			}
			if importName != tt.expectedImport {
				t.Fatalf("expected import name to be %q, got %q", tt.expectedImport, importName)
			}
		})
	}
}

// This loads a file, and determines which functions are exported
// and match a GoGo function signature.
func TestFindGoGoFunctions(t *testing.T) {
	methodNames := []string{
		"Description",
		"SetHelp",
		"Argument",
		"DescriptionArgument",
		"CtxChained",
		"ArgumentChained",
	}
	tests := []struct {
		name      string
		prefix    string
		gogoAlias string
		filePath  string
		expected  any
		err       bool
	}{
		{
			name:      "basic",
			prefix:    "Basic",
			gogoAlias: "gogo",
			filePath:  "testdata/basic.go",
		},
		{
			name:      "aliased",
			prefix:    "Aliased",
			gogoAlias: "gogo2",
			filePath:  "testdata/aliased.go",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse the source code
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, tt.filePath, nil, 0)
			if err != nil {
				t.Fatalf("error parsing file: %v", err)
			}
			funcs := getGoGoFunctions(tt.gogoAlias, file)
			prefixedNames := make([]string, len(methodNames))
			for i, name := range methodNames {
				prefixedNames[i] = tt.prefix + name
			}
			if !validateFunctions(prefixedNames, funcs) {
				t.Fatalf("expected functions to match %v", prefixedNames)
			}
		})
	}
}

// Tests to insure we can find the GoGo Context and Argument options
// and their values, whether chained or not.

//func TestFindGoGoOptions(t *testing.T) {
//	methodNames := []string{
//		"Description",
//		"SetHelp",
//		"Argument",
//		"DescriptionArgument",
//		"CtxChained",
//		"ArgumentChained",
//	}
//	tests := []struct {
//		name     string
//		funcName   string
//		filePath string
//		expected any
//		err      bool
//	}{
//		{
//			name:     "basic description",
//			funcName:   "Basic",
//			filePath: "testdata/basic.go",
//			expected: "set a description",
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			// load the file
//			f, err := os.ReadFile(tt.filePath)
//			if err != nil {
//				t.Fatalf("error reading file: %v", err)
//			}
//			if f == nil {
//				t.Fatal("expected file to be read")
//			}
//			// for each methodName
//		})
//	}
//}

// Test that all the various options are found
// when they are not chained together.
func TestNoChainedOptions(t *testing.T) {
	tests := []struct {
		name       string
		funcName   string
		optionName string
		expected   any
		err        bool
	}{
		{
			name:       "ctx description",
			funcName:   "BasicDescription",
			optionName: "SetDescription",
			expected:   "set a description",
		},
		{
			name:       "ctx set help",
			funcName:   "BasicSetHelp",
			optionName: "SetHelp",
			expected:   "is this a thing?",
		},
	}
	// Parse the source code
	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, "testdata/basic.go", nil, 0)
	if err != nil {
		t.Fatalf("error parsing file: %v", err)
	}
	alias, err := getGoGoImportName(astFile)
	if err != nil {
		t.Fatalf("error getting import name: %v", err)
	}
	funcs := getGoGoFunctions(alias, astFile)
	if funcs == nil {
		t.Fatal("expected functions to be found")
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// for each funcNames, retrieve that function
			funk := getFuncByName(funcs, tt.funcName)
			if funk == nil {
				t.Fatalf("expected function %q to be found", tt.funcName)
			}
		})
	}
}

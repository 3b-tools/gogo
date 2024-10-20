package gogo

import (
	"embed"
	"github.com/bitfield/script"
	"os"
	"strings"
	"text/template"
)

/*
thoughts on how this works:
1. load a file with functions
2. parse the file into functions
3. generate the binary with


3. find the function that matches the requested function
4. parse arguments from the cmdLine or input into the desired types
5. pass the arguments to the function

Where does compilation come into play?
How can I debug the process?
*/

//go:embed templates/*
var templates embed.FS

type renderData struct {
	RootCmd     GoCmd
	SubCommands []GoCmd
}

type GoCmd struct {
	Name    string // Name of the command
	Short   string // Short Description of the command
	Long    string // Long Description of the command
	Body    string // Body of the command
	GoFlags []GoFlag
}

type GoFlag struct {
	Type    string // string, int, bool, float64
	Name    string // name of the flag
	Short   string // short name of the flag
	Default any    // default value of the flag
	Help    string // help text for the flag
}

func RenderTemplates(rd renderData) (string, error) {
	tmpl, err := template.ParseFS(templates,
		"templates/main.go.tmpl",
		"templates/subCmd.go.tmpl",
	)
	if err != nil {
		return "", err
	}
	// make a string buffer to write to
	outBuf := new(strings.Builder)
	err = tmpl.Execute(outBuf, rd)
	if err != nil {
		return "", err
	}
	return outBuf.String(), nil
}

// writeFile writes the input to the location
func writeFile(location, input string) error {
	return os.WriteFile(location, []byte(input), 0644)
}

// go fmt
func runGoFmt(location string) error {
	_, err := script.Exec("go fmt " + location).String()
	return err
}

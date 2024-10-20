package gogo

import (
	"embed"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/3b-tools/gogo/run"

	"github.com/bitfield/script"
)

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

func Build(log *log.Logger) {

}

func RenderTemplates(rd renderData) (string, error) {
	tmpl, err := template.ParseFS(templates,
		"templates/main.go.tmpl",
		"templates/subCmd.go.tmpl",
	)
	if err != nil {
		return "", err
	}
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
func runGoFmt(log *log.Logger, location string) error {
	_, err := script.Exec("go fmt " + location).String()
	return err
}

func buildBinary(log *log.Logger, dir string) error {
	//_ = run.CMD("go", "get").Dir(dir).Log().Run()
	//err = run.CMD("go", "build").Dir(dir).Log().Run()
	//if err != nil {
	//	os.Remove(filepath.Join(dir, "bake"))
	//	return fmt.Errorf("failed to build binary: %w", err)
	//}
	runLine := fmt.Sprintf("go get -d %s", dir)
	//go get
	getOutput, err := script.Exec(runLine).String()
	if err != nil {
		return fmt.Errorf("failed to get dependencies: %s : %w", getOutput, err)
	}
	// build
	runLine = fmt.Sprintf("go build %s", dir)
	return run.CMD(log, "go", "build").Dir(dir).Log().Run()
}

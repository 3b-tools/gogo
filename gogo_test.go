package gogo

import (
	"os"
	"testing"
)

func TestRenderTemplates(t *testing.T) {
	tests := []struct {
		name       string
		renderData renderData
	}{
		{
			name: "empty",
			renderData: renderData{
				RootCmd: GoCmd{
					Name:    "rootFlag",
					GoFlags: nil,
				},
			},
		},
		{
			name: "rootcmd with flags",
			renderData: renderData{
				RootCmd: GoCmd{
					Name: "rootFlag",
					GoFlags: []GoFlag{
						{
							Type:    "string",
							Name:    "stringFlag",
							Short:   "s",
							Default: "default",
							Help:    "help text",
						},
					},
				},
			},
		},
		{
			name: "subCmd with flags",
			renderData: renderData{
				RootCmd: GoCmd{
					Name:    "rootFlag",
					GoFlags: nil,
				},
				SubCommands: []GoCmd{
					{
						Name: "subCmd",
						GoFlags: []GoFlag{
							{
								Type:    "string",
								Name:    "stringFlag",
								Short:   "s",
								Default: "default",
								Help:    "help text",
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := RenderTemplates(tt.renderData)
			if err != nil {
				t.Fatal(err)
			}
			t.Log(res)
		})
	}
}

func TestRenderGoFmt(t *testing.T) {
	data := renderData{
		RootCmd: GoCmd{},
		SubCommands: []GoCmd{
			{
				Name:  "subCmd",
				Short: "short description",
				Long:  "long description",
				GoFlags: []GoFlag{
					{
						Type:    "string",
						Name:    "stringFlag",
						Short:   "s",
						Default: "default",
					},
				},
			},
		},
	}
	res, err := RenderTemplates(data)
	if err != nil {
		t.Fatal(err)
	}
	// now writeFile and run gofmt
	target := "/tmp/out.go"
	err = writeFile(target, res)
	if err != nil {
		t.Fatal(err)
	}
	err = runGoFmt(nil, target)
	if err != nil {
		t.Fatal(err)
	}
	// read file and output it
	formattedRes, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(formattedRes))
}

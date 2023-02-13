package main

import (
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"go/format"
	"os"
	"path"
	"path/filepath"
	"text/template"

	"github.com/go-leo/stringx"
)

var moduleName = flag.String("module", "", "module name")

var appPath = flag.String("app", "", "app path")

//go:embed cmd_wire.go.template
var cmdWireContent string

//go:embed app_root_wire.go.template
var appRootWireContent string

//go:embed presentation_wire.go.template
var presentationWireContent string

//go:embed sample_wire.go.template
var sampleWireContent string

//go:embed bus_commands.go.template
var busCommandsContent string

//go:embed bus_queries.go.template
var busQueriesContent string

//go:embed bus_wire.go.template
var busContent string

//go:embed application_wire.go.template
var applicationContent string

//go:embed infrastructure_wire.go.template
var infrastructureContent string

func main() {
	flag.Parse()
	if stringx.IsBlank(*moduleName) {
		fmt.Println("module name is empty")
		flag.Usage()
		return
	}
	if stringx.IsBlank(*appPath) {
		fmt.Println("app path is empty")
		flag.Usage()
		return
	}

	sources := []*Source{
		newSource(path.Join("cmd", *appPath), cmdWireContent, "wire.go"),
		newSource(path.Join("internal/app", *appPath), appRootWireContent, "wire.go"),

		newSource(path.Join("internal/app", *appPath, "config"), sampleWireContent, "wire.go"),

		newSource(path.Join("internal/app", *appPath, "presentation"), presentationWireContent, "wire.go"),
		newSource(path.Join("internal/app", *appPath, "presentation", "bus"), busContent, "wire.go"),
		newSource(path.Join("internal/app", *appPath, "presentation", "bus"), busCommandsContent, "commands.go"),
		newSource(path.Join("internal/app", *appPath, "presentation", "bus"), busQueriesContent, "queries.go"),
		newSource(path.Join("internal/app", *appPath, "presentation", "assembler"), sampleWireContent, "wire.go"),
		newSource(path.Join("internal/app", *appPath, "presentation", "console"), sampleWireContent, "wire.go"),
		newSource(path.Join("internal/app", *appPath, "presentation", "controller"), sampleWireContent, "wire.go"),
		newSource(path.Join("internal/app", *appPath, "presentation", "provider"), sampleWireContent, "wire.go"),
		newSource(path.Join("internal/app", *appPath, "presentation", "resource"), sampleWireContent, "wire.go"),
		newSource(path.Join("internal/app", *appPath, "presentation", "subscriber"), sampleWireContent, "wire.go"),

		newSource(path.Join("internal/app", *appPath, "application"), applicationContent, "wire.go"),
		newSource(path.Join("internal/app", *appPath, "application", "command"), sampleWireContent, "wire.go"),
		newSource(path.Join("internal/app", *appPath, "application", "query"), sampleWireContent, "wire.go"),

		newSource(path.Join("internal/app", *appPath, "domain"), sampleWireContent, "wire.go"),

		newSource(path.Join("internal/app", *appPath, "infrastructure"), infrastructureContent, "wire.go"),
		newSource(path.Join("internal/app", *appPath, "infrastructure", "clientadapter"), sampleWireContent, "wire.go"),
		newSource(path.Join("internal/app", *appPath, "infrastructure", "clientport"), sampleWireContent, "wire.go"),
		newSource(path.Join("internal/app", *appPath, "infrastructure", "converter"), sampleWireContent, "wire.go"),
		newSource(path.Join("internal/app", *appPath, "infrastructure", "publisheradapter"), sampleWireContent, "wire.go"),
		newSource(path.Join("internal/app", *appPath, "infrastructure", "publisherport"), sampleWireContent, "wire.go"),
		newSource(path.Join("internal/app", *appPath, "infrastructure", "repositoryadapter"), sampleWireContent, "wire.go"),
		newSource(path.Join("internal/app", *appPath, "infrastructure", "repositoryport"), sampleWireContent, "wire.go"),
	}
	for _, src := range sources {
		err := createSource(src)
		if err != nil {
			panic(err)
		}
	}

}

func newSource(dirPath string, text string, name string) *Source {
	wireFilePath := filepath.Join(dirPath, name)
	src := &Source{
		dirPath:  dirPath,
		filePath: wireFilePath,
		text:     text,
		data: &SourceData{
			ModuleName:   *moduleName,
			AppPath:      *appPath,
			AppBaseName:  filepath.Base(*appPath),
			WireFilePath: wireFilePath,
			Package:      filepath.Base(dirPath),
		},
	}
	return src
}

type SourceData struct {
	ModuleName   string
	AppPath      string
	AppBaseName  string
	WireFilePath string
	Package      string
}

type Source struct {
	dirPath  string
	filePath string
	text     string
	data     *SourceData
}

func createSource(src *Source) error {
	err := os.MkdirAll(src.dirPath, 0777)
	if err != nil {
		return err
	}
	tmpl, err := template.New(src.filePath).Parse(src.text)
	if err != nil {
		return err
	}
	file, err := os.Create(src.filePath)
	if err != nil {
		return err
	}
	buffer := &bytes.Buffer{}
	err = tmpl.Execute(buffer, src.data)
	if err != nil {
		return err
	}
	source, err := format.Source(buffer.Bytes())
	if err != nil {
		return err
	}
	_, err = file.Write(source)
	if err != nil {
		return err
	}
	return nil
}

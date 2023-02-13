package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"os"
	"path"
	"text/template"

	"github.com/go-leo/stringx"
)

var moduleName = flag.String("module", "", "module name")

var appPath = flag.String("app", "", "app path")

type packageInfo struct {
	moduleName string
	appPath    string
}

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
	fmt.Println(*moduleName)
	fmt.Println(*appPath)
	info := &packageInfo{
		moduleName: *moduleName,
		appPath:    *appPath,
	}
	err := makeDir(info)
	if err != nil {
		panic(err)
	}
}

type WireFile struct {
	Path      string
	IsMain    bool
	Package   string
	Imports   []string
	Providers []string
}

func makeDir(info *packageInfo) error {
	var wireFiles []*WireFile

	appRootPath := path.Join("internal/app", info.appPath)
	wireFiles = appendCmd(wireFiles, info, appRootPath)
	wireFiles = appendApp(wireFiles, info, appRootPath)

	for _, files := range wireFiles {
		err := os.MkdirAll(files.Path, 0777)
		if err != nil {
			return err
		}
		tmpl, err := template.New("test").Parse(wireFile)
		if err != nil {
			return err
		}
		file, err := os.Create(path.Join(files.Path, "wire.go"))
		if err != nil {
			return err
		}
		buffer := &bytes.Buffer{}
		err = tmpl.Execute(buffer, files)
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
	}
	return nil
}

func appendCmd(wireFiles []*WireFile, info *packageInfo, appRootPath string) []*WireFile {
	cmdPath := path.Join("cmd", info.appPath)
	wireFiles = append(wireFiles, &WireFile{
		Path:    cmdPath,
		IsMain:  true,
		Package: "main",
		Imports: []string{
			path.Join(info.moduleName, appRootPath),
		},
		Providers: []string{
			path.Base(appRootPath),
		},
	})
	return wireFiles
}

func appendApp(wireFiles []*WireFile, info *packageInfo, appRootPath string) []*WireFile {
	configPath := path.Join("internal/app", info.appPath, "config")
	presentationPath := path.Join("internal/app", info.appPath, "presentation")
	applicationPath := path.Join("internal/app", info.appPath, "application")
	domainPath := path.Join("internal/app", info.appPath, "domain")
	infrastructurePath := path.Join("internal/app", info.appPath, "infrastructure")
	wireFiles = appendRoot(wireFiles, info, appRootPath, configPath, presentationPath, applicationPath, domainPath, infrastructurePath)
	wireFiles = appendConfig(wireFiles, configPath)
	wireFiles = appendPresentation(wireFiles, info, presentationPath)
	wireFiles = appendApplication(wireFiles, info, applicationPath)
	wireFiles = appendDomain(wireFiles, info, domainPath)
	wireFiles = appendInfrastructure(wireFiles, info, infrastructurePath)
	return wireFiles
}

func appendRoot(wireFiles []*WireFile, info *packageInfo, appRootPath string, configPath string, presentationPath string, applicationPath string, domainPath string, infrastructurePath string) []*WireFile {
	appName := path.Base(info.appPath)
	wireFiles = append(wireFiles, &WireFile{
		Path:    appRootPath,
		IsMain:  false,
		Package: appName,
		Imports: []string{
			path.Join(info.moduleName, configPath),
			path.Join(info.moduleName, presentationPath),
			path.Join(info.moduleName, applicationPath),
			path.Join(info.moduleName, domainPath),
			path.Join(info.moduleName, infrastructurePath),
		},
		Providers: []string{
			path.Base(configPath),
			path.Base(applicationPath),
			path.Base(domainPath),
			path.Base(infrastructurePath),
			path.Base(presentationPath),
		},
	})
	return wireFiles
}

func appendDomain(wireFiles []*WireFile, info *packageInfo, domainPath string) []*WireFile {
	internalPath := path.Join("internal/app", info.appPath, "domain/internal")
	domainnamePath := path.Join("internal/app", info.appPath, "domain/domainname")
	wireFiles = append(wireFiles, &WireFile{
		Path:    domainPath,
		IsMain:  false,
		Package: "domain",
		Imports: []string{
			path.Join(info.moduleName, internalPath),
			path.Join(info.moduleName, domainnamePath),
		},
		Providers: []string{
			path.Base(internalPath),
			path.Base(domainnamePath),
		},
	})
	wireFiles = append(wireFiles, &WireFile{
		Path:    internalPath,
		IsMain:  false,
		Package: "internal"})
	wireFiles = append(wireFiles, &WireFile{
		Path:    domainnamePath,
		IsMain:  false,
		Package: "domainname"})
	return wireFiles
}

func appendConfig(wireFiles []*WireFile, configPath string) []*WireFile {
	wireFiles = append(wireFiles, &WireFile{Path: configPath, IsMain: false, Package: "config"})
	return wireFiles
}

func appendApplication(wireFiles []*WireFile, info *packageInfo, applicationPath string) []*WireFile {
	commandPath := path.Join("internal/app", info.appPath, "application/command")
	queryPath := path.Join("internal/app", info.appPath, "application/query")
	servicePath := path.Join("internal/app", info.appPath, "application/service")
	wireFiles = append(wireFiles, &WireFile{
		Path:    applicationPath,
		IsMain:  false,
		Package: "application",
		Imports: []string{
			path.Join(info.moduleName, commandPath),
			path.Join(info.moduleName, queryPath),
			path.Join(info.moduleName, servicePath),
		},
		Providers: []string{
			path.Base(commandPath),
			path.Base(queryPath),
			path.Base(servicePath),
		},
	})
	wireFiles = append(wireFiles, &WireFile{
		Path:    commandPath,
		IsMain:  false,
		Package: "command",
	})
	wireFiles = append(wireFiles, &WireFile{
		Path:    queryPath,
		IsMain:  false,
		Package: "query",
	})
	wireFiles = append(wireFiles, &WireFile{
		Path:    servicePath,
		IsMain:  false,
		Package: "service",
	})
	return wireFiles
}

func appendInfrastructure(wireFiles []*WireFile, info *packageInfo, appInfrastructurePath string) []*WireFile {
	adapterPath := path.Join("internal/app", info.appPath, "infrastructure/adapter")
	portPath := path.Join("internal/app", info.appPath, "infrastructure/port")
	converterPath := path.Join("internal/app", info.appPath, "infrastructure/converter")
	wireFiles = append(wireFiles, &WireFile{
		Path:    appInfrastructurePath,
		IsMain:  false,
		Package: "infrastructure",
		Imports: []string{
			path.Join(info.moduleName, adapterPath),
			path.Join(info.moduleName, portPath),
			path.Join(info.moduleName, converterPath),
		},
		Providers: []string{
			path.Base(adapterPath),
			path.Base(portPath),
			path.Base(converterPath),
		},
	})

	wireFiles = append(wireFiles, &WireFile{
		Path:    converterPath,
		IsMain:  false,
		Package: "converter",
	})

	adapterClientPath := path.Join("internal/app", info.appPath, "infrastructure/adapter/client")
	adapterPublisherPath := path.Join("internal/app", info.appPath, "infrastructure/adapter/publisher")
	adapterRepositoryPath := path.Join("internal/app", info.appPath, "infrastructure/adapter/repository")
	wireFiles = append(wireFiles, &WireFile{
		Path:    adapterPath,
		IsMain:  false,
		Package: "adapter",
		Imports: []string{
			path.Join(info.moduleName, adapterClientPath),
			path.Join(info.moduleName, adapterPublisherPath),
			path.Join(info.moduleName, adapterRepositoryPath),
		},
		Providers: []string{
			path.Base(adapterClientPath),
			path.Base(adapterPublisherPath),
			path.Base(adapterRepositoryPath),
		},
	})
	portClientPath := path.Join("internal/app", info.appPath, "infrastructure/port/client")
	portPublisherPath := path.Join("internal/app", info.appPath, "infrastructure/port/publisher")
	portRepositoryPath := path.Join("internal/app", info.appPath, "infrastructure/port/repository")
	wireFiles = append(wireFiles, &WireFile{
		Path:    portPath,
		IsMain:  false,
		Package: "port",
		Imports: []string{
			path.Join(info.moduleName, portClientPath),
			path.Join(info.moduleName, portPublisherPath),
			path.Join(info.moduleName, portRepositoryPath),
		},
		Providers: []string{
			path.Base(portClientPath),
			path.Base(portPublisherPath),
			path.Base(portRepositoryPath),
		},
	})

	wireFiles = append(wireFiles, &WireFile{
		Path:    adapterClientPath,
		IsMain:  false,
		Package: "client",
	})
	wireFiles = append(wireFiles, &WireFile{
		Path:    adapterPublisherPath,
		IsMain:  false,
		Package: "publisher",
	})
	wireFiles = append(wireFiles, &WireFile{
		Path:    adapterRepositoryPath,
		IsMain:  false,
		Package: "repository",
	})

	wireFiles = append(wireFiles, &WireFile{
		Path:    portClientPath,
		IsMain:  false,
		Package: "client"})
	wireFiles = append(wireFiles, &WireFile{
		Path:    portPublisherPath,
		IsMain:  false,
		Package: "publisher"})
	wireFiles = append(wireFiles, &WireFile{
		Path:    portRepositoryPath,
		IsMain:  false,
		Package: "repository"})
	return wireFiles
}

func appendPresentation(wireFiles []*WireFile, info *packageInfo, appPresentationPath string) []*WireFile {
	presentationAdapterPath := path.Join("internal/app", info.appPath, "presentation/adapter")
	presentationPortPath := path.Join("internal/app", info.appPath, "presentation/port")
	presentationBus := path.Join("internal/app", info.appPath, "presentation/bus")
	presentationAssembler := path.Join("internal/app", info.appPath, "presentation/assembler")
	wireFiles = append(wireFiles, &WireFile{
		Path:    appPresentationPath,
		IsMain:  false,
		Package: "presentation",
		Imports: []string{
			path.Join(info.moduleName, presentationAdapterPath),
			path.Join(info.moduleName, presentationPortPath),
			path.Join(info.moduleName, presentationBus),
			path.Join(info.moduleName, presentationAssembler),
		},
		Providers: []string{
			path.Base(presentationAdapterPath),
			path.Base(presentationPortPath),
			path.Base(presentationBus),
			path.Base(presentationAssembler),
		},
	})

	adapterConsole := path.Join("internal/app", info.appPath, "presentation/adapter/console")
	adapterController := path.Join("internal/app", info.appPath, "presentation/adapter/controller")
	adapterResource := path.Join("internal/app", info.appPath, "presentation/adapter/resource")
	adapterProvider := path.Join("internal/app", info.appPath, "presentation/adapter/provider")
	adapterSubscriber := path.Join("internal/app", info.appPath, "presentation/adapter/subscriber")
	wireFiles = append(wireFiles, &WireFile{
		Path:    presentationAdapterPath,
		IsMain:  false,
		Package: "adapter",
		Imports: []string{
			path.Join(info.moduleName, adapterConsole),
			path.Join(info.moduleName, adapterController),
			path.Join(info.moduleName, adapterResource),
			path.Join(info.moduleName, adapterProvider),
			path.Join(info.moduleName, adapterSubscriber),
		},
		Providers: []string{
			path.Base(adapterConsole),
			path.Base(adapterController),
			path.Base(adapterResource),
			path.Base(adapterProvider),
			path.Base(adapterSubscriber),
		},
	})
	wireFiles = append(wireFiles, &WireFile{
		Path:    presentationPortPath,
		IsMain:  false,
		Package: "port"})
	wireFiles = append(wireFiles, &WireFile{
		Path:    presentationAssembler,
		IsMain:  false,
		Package: "assembler"})
	wireFiles = append(wireFiles, &WireFile{
		Path:    presentationBus,
		IsMain:  false,
		Package: "bus"})

	wireFiles = append(wireFiles, &WireFile{
		Path:    adapterConsole,
		IsMain:  false,
		Package: "console"})
	wireFiles = append(wireFiles, &WireFile{
		Path:    adapterController,
		IsMain:  false,
		Package: "controller"})
	wireFiles = append(wireFiles, &WireFile{
		Path:    adapterResource,
		IsMain:  false,
		Package: "resource"})
	wireFiles = append(wireFiles, &WireFile{
		Path:    adapterProvider,
		IsMain:  false,
		Package: "provider"})
	wireFiles = append(wireFiles, &WireFile{
		Path:    adapterSubscriber,
		IsMain:  false,
		Package: "subscriber"})
	return wireFiles
}

var wireFile = `{{- if .IsMain }}
//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.
{{- end}}
package {{.Package}}

import (
	"github.com/google/wire"
{{- range .Imports }}
"{{ . }}"
{{- end}}
)

var Provider = wire.NewSet(
{{- range .Providers }}
{{ . }}.Provider,
{{- end}}
)
`

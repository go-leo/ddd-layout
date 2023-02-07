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
	appName := path.Base(info.appPath)
	appRootPath := path.Join("internal/app", info.appPath)

	appConfigPath := path.Join("internal/app", info.appPath, "config")
	appCorePath := path.Join("internal/app", info.appPath, "core")
	appInfrastructurePath := path.Join("internal/app", info.appPath, "infrastructure")
	appPresentationPath := path.Join("internal/app", info.appPath, "presentation")

	applicationPath := path.Join("internal/app", info.appPath, "core/application")
	domainPath := path.Join("internal/app", info.appPath, "core/domain")

	commandPath := path.Join("internal/app", info.appPath, "core/application/command")
	queryPath := path.Join("internal/app", info.appPath, "core/application/query")
	servicePath := path.Join("internal/app", info.appPath, "core/application/service")

	internalPath := path.Join("internal/app", info.appPath, "core/domain/internal")
	domainnamePath := path.Join("internal/app", info.appPath, "core/domain/domainname")

	var wireFiles []*WireFile
	wireFiles = append(wireFiles, &WireFile{
		Path:    path.Join("cmd", info.appPath),
		IsMain:  true,
		Package: "main",
		Imports: []string{
			path.Join(info.moduleName, appRootPath),
		},
		Providers: []string{
			path.Base(appRootPath),
		},
	})

	wireFiles = append(wireFiles, &WireFile{
		Path:    appRootPath,
		IsMain:  false,
		Package: appName,
		Imports: []string{
			path.Join(info.moduleName, appConfigPath),
			path.Join(info.moduleName, appCorePath),
			path.Join(info.moduleName, appInfrastructurePath),
			path.Join(info.moduleName, appPresentationPath),
		},
		Providers: []string{
			path.Base(appConfigPath),
			path.Base(appCorePath),
			path.Base(appInfrastructurePath),
			path.Base(appPresentationPath),
		},
	})

	wireFiles = append(wireFiles, &WireFile{
		Path:    appConfigPath,
		IsMain:  false,
		Package: "config"})

	wireFiles = append(wireFiles, &WireFile{
		Path:    appCorePath,
		IsMain:  false,
		Package: "core",
		Imports: []string{
			path.Join(info.moduleName, applicationPath),
			path.Join(info.moduleName, domainPath),
		},
		Providers: []string{
			path.Base(applicationPath),
			path.Base(domainPath),
		},
	})

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

	wireFiles = append(wireFiles, &WireFile{
		Path:    internalPath,
		IsMain:  false,
		Package: "internal"})
	wireFiles = append(wireFiles, &WireFile{
		Path:    domainnamePath,
		IsMain:  false,
		Package: "domainname"})

	adapterPath := path.Join("internal/app", info.appPath, "infrastructure/adapter")
	portPath := path.Join("internal/app", info.appPath, "infrastructure/port")
	wireFiles = append(wireFiles, &WireFile{
		Path:    appInfrastructurePath,
		IsMain:  false,
		Package: "infrastructure",
		Imports: []string{
			path.Join(info.moduleName, adapterPath),
			path.Join(info.moduleName, portPath),
		},
		Providers: []string{
			path.Base(adapterPath),
			path.Base(portPath),
		},
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

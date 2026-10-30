package templates

import (
	"embed"

	"github.com/nmusey/letsgo/lib/files"
	"github.com/nmusey/letsgo/lib/names"
)

//go:embed all:project
var projectFilesystem embed.FS

func NewProjectTemplate(name string, repository string, dbPort string) *files.TemplateParser {
	return &files.TemplateParser{
		Filesystem:     projectFilesystem,
		FilesystemRoot: "project",
		Substitutions: map[string]string{
			"AppName": name,
			"Module":  repository,
			"DbPort": dbPort,
		},
	}
}

//go:embed all:model
var modelFilesystem embed.FS

func NewModelTemplate(name string) (*files.TemplateParser, error) {
	moduleName, err := files.ModuleName()
	if err != nil {
		return nil, err
	}

	return &files.TemplateParser{
		Filesystem:     modelFilesystem,
		FilesystemRoot: "model",
		Substitutions: map[string]string{
			"Package": names.PackageName(name),
			"Struct":  names.StructName(name),
			"Module":  moduleName,
		},
	}, nil
}

//go:embed all:appservice
var appserviceFilesystem embed.FS

func NewAppServiceTemplate(name string) (*files.TemplateParser, error) {
	moduleName, err := files.ModuleName()
	if err != nil {
		return nil, err
	}

	return &files.TemplateParser{
		Filesystem:     appserviceFilesystem,
		FilesystemRoot: "appservice",
		Substitutions: map[string]string{
			"Struct":  names.StructName(name),
			"Package": names.PackageName(name),
			"Module":  moduleName,
		},
	}, nil
}

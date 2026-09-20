package templates

import (
	"embed"

	"github.com/nmusey/letsgo/lib/files"
)

//go:embed all:project
var projectFilesystem embed.FS

func NewProjectTemplate(name string, repository string) *files.TemplateParser {
	return &files.TemplateParser{
		Filesystem:     projectFilesystem,
		FilesystemRoot: "project",
		Substitutions: map[string]string{
			"AppName": name,
			"AppRepo": repository,
		},
	}
}

//go:embed all:model
var modelFilesystem embed.FS

func NewModelTemplate(name string) *files.TemplateParser {
	return &files.TemplateParser{
		Filesystem:     modelFilesystem,
		FilesystemRoot: "model",
		Substitutions: map[string]string{
			"Package": "", // TODO
			"Struct":  "", // TODO
			"AppName": "", // TODO
			"AppRepo": "", // TODO
		},
	}
}

//go:embed all:appservice
var appserviceFilesystem embed.FS

func NewAppServiceTemplate(name string) *files.TemplateParser {
	return &files.TemplateParser{
		Filesystem:     appserviceFilesystem,
		FilesystemRoot: "appservice",
		Substitutions: map[string]string{
			"ServiceName": name,
			"AppName":     "", // TODO
			"AppRepo":     "", // TODO
		},
	}
}

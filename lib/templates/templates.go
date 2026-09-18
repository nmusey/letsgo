package templates

import (
	"embed"

	"github.com/nmusey/letsgo/lib/files"
)

//go:embed all:project
var projectFilesystem embed.FS

func NewProjectTemplate(name string, repository string) *files.TemplateParser {
	return &files.TemplateParser{
		Path: name,
		Filesystem:     projectFilesystem,
		FilesystemRoot: "project",
		Substitutions: map[string]string{
			"AppName": name,
			"AppRepo": repository,
		},
	}
}

//go:embed all:generators
var generatorsFilesystem embed.FS

func GeneratorsFilesystem() embed.FS {
	return generatorsFilesystem
}

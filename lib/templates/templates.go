package templates

import (
	"embed"

	"github.com/nmusey/letsgo/lib/files"
)

//go:embed all:project
var projectFilesystem embed.FS

func NewProjectTemplate(name string) *files.TemplateParser {
	return &files.TemplateParser{
		Path:       name,
		Filesystem: projectFilesystem,
		Substitutions: map[string]string{
			"AppName": name,
		},
	}
}

//go:embed all:test
var testFilesystem embed.FS

// Template for internal testing only
func NewTestTemplate(text string) *files.TemplateParser {
	return &files.TemplateParser{
		Path: "test",
		Filesystem: testFilesystem,
		Substitutions: map[string]string{
			"Test": text,
		},
	}
}

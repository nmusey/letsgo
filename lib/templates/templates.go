package templates

import "embed"

type Template struct {
	Filesystem embed.FS
	Path       string
}

//go:embed all:project
var projectFilesystem embed.FS
var ProjectTemplate = Template{
	Filesystem: projectFilesystem,
	Path: "project",
}

package main

import (
	"io"
	"io/fs"
	"os"
	"text/template"

	"github.com/nmusey/letsgo/lib/templates"
)

// TODO this should be more flexible
type Substitutions struct {
	AppName string
	DbPort  int
}

func main() {
	substitutions := Substitutions{
		AppName: "App Name",
		DbPort:  5432,
	}

	err := parseFiles(templates.ProjectTemplate, substitutions)
	if err != nil {
		panic(err)
	}
}

func parseFiles(template templates.Template, substitutions Substitutions) error {
	err := fs.WalkDir(template.Filesystem, template.Path, func(path string, entry fs.DirEntry, err error) error {
		if entry.IsDir() {
			return nil
		}

		targetPrefix := "tmp/"
		os.MkdirAll(targetPrefix, os.ModePerm)
		file, err := os.Create(targetPrefix+entry.Name())
		if err != nil {
			return err
		}

		defer file.Close()

		pathPrefix := "lib/templates/"
		err = substituteFile(pathPrefix+path, file, substitutions)
		return err
	})

	return err
}

func substituteFile(filename string, writer io.Writer, substitutions Substitutions) error {
	contents, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	tmpl, err := template.New(filename).Parse(string(contents))
	if err != nil {
		return err
	}

	err = tmpl.Execute(writer, substitutions)
	if err != nil {
		return err
	}

	return nil
}

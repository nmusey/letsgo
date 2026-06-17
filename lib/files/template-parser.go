package files

import (
	"embed"
	"html/template"
	"io"
	"io/fs"
	"os"
)

type TemplateParser struct {
	Path       string
	Filesystem embed.FS

	Substitutions map[string]string
}

func (t TemplateParser) ParseFiles() error {
	return fs.WalkDir(t.Filesystem, ".", func(entryPath string, entry fs.DirEntry, err error) error {
		if entry == nil {
			return nil
		}

		if entry.IsDir() {
			os.Mkdir(entryPath, 0777)
			return nil
		}

		file, err := os.Create(entryPath)
		if err != nil {
			return err
		}

		defer file.Close()

		outpath := entryPath
		err = t.SubstituteFile(outpath, file)

		return err
	})
}

func (t TemplateParser) SubstituteFile(filename string, writer io.Writer) error {
	contents, err := fs.ReadFile(t.Filesystem, filename)
	if err != nil {
		return err
	}

	tmpl, err := template.New(filename).Parse(string(contents))
	if err != nil {
		return err
	}

	err = tmpl.Execute(writer, t.Substitutions)
	if err != nil {
		return err
	}

	return nil
}

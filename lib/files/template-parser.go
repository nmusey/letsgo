package files

import (
	"embed"
	"html/template"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type TemplateParser struct {
	Path       string

	Filesystem embed.FS
	FilesystemRoot string

	Substitutions map[string]string
}

func (t TemplateParser) ParseFiles() error {
	os.Mkdir(t.Path, 0755)

	return fs.WalkDir(t.Filesystem, t.FilesystemRoot, func(entryPath string, entry fs.DirEntry, err error) error {
		if entry == nil {
			return nil
		}

		outpath := t.transformToOutputPath(entryPath)
		if entry.IsDir() {
			os.Mkdir(outpath, 0777) // TODO - can this be 0755?
			return nil
		}

		file, err := os.Create(outpath)
		if err != nil {
			return err
		}

		defer file.Close()
		return t.SubstituteFile(entryPath, file)
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

// TODO - clean up the logic here and in calling locations
func RenderTemplate(fsys fs.FS, templatePath, outPath string, substitutions map[string]string) error {
	contents, err := fs.ReadFile(fsys, templatePath)
	if err != nil {
		return err
	}

	tmpl, err := template.New(templatePath).Parse(string(contents))
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(outPath), 0777); err != nil {
		return err
	}

	file, err := os.OpenFile(outPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, substitutions)
}

// TODO - function documentation pass https://go.dev/doc/comment
func (t TemplateParser) transformToOutputPath(entryPath string) string {
		outpath := strings.Replace(entryPath, t.FilesystemRoot, t.Path, 1)
		return strings.TrimSuffix(outpath, ".template")
}

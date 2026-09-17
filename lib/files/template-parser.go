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

		outpath, templateFileFound := strings.CutSuffix(entryPath, ".template")
		file, err := os.Create(outpath)
		if err != nil {
			return err
		}

		defer file.Close()

		if templateFileFound {
			return t.SubstituteFile(entryPath, file)
		} 
		
		return t.CopyFile(entryPath, file)
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

func (t TemplateParser) CopyFile(filename string, writer io.Writer) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	defer file.Close()
	_, err = io.Copy(writer, file)

	return err
}

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

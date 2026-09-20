package files

import (
	"embed"
	"io"
	"io/fs"
	"os"
	"strings"
	"text/template"
)

// Holds information needed to parse templates
type TemplateParser struct {
	Filesystem     embed.FS
	FilesystemRoot string
	OutputPath     *string

	Substitutions map[string]string
}

// Parse the files in a given template to a new directory with given name
func (t TemplateParser) ParseFiles(dirname string) error {
	os.MkdirAll(dirname, 0755)

	return fs.WalkDir(t.Filesystem, t.FilesystemRoot, func(entryPath string, entry fs.DirEntry, err error) error {
		if entry == nil {
			return nil
		}

		outpath := t.getOutpath(entryPath, dirname)
		if entry.IsDir() {
			os.Mkdir(outpath, 0755)
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

// Reads a file and substitute and variables in the parser
func (t TemplateParser) SubstituteFile(filename string, writer io.Writer) error {
	// The entire file is read at once, template files are relatively small so this isn't an issue.
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

// Turn a path from the read filesystem to the write location.
func (t TemplateParser) getOutpath(entryPath string, outputDir string) string {
		outpath := entryPath
		outpath = strings.Replace(outpath, t.FilesystemRoot, outputDir, 1)
		outpath = strings.TrimSuffix(outpath, ".template")

		return outpath
}

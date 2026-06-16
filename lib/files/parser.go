package files

import (
	"io/fs"
	"os"

	"github.com/nmusey/letsgo/lib/templates"
)

func ParseFiles(template templates.Template, substitutions Substitutions) error {
	err := fs.WalkDir(template.Filesystem, template.Path, func(path string, entry fs.DirEntry, err error) error {
		if entry.IsDir() {
			return nil
		}

		targetPrefix := "tmp/" // TODO This should be configurable
		os.MkdirAll(targetPrefix, os.ModePerm)
		file, err := os.Create(targetPrefix+entry.Name())
		if err != nil {
			return err
		}

		defer file.Close()

		pathPrefix := "lib/templates/" // TODO this should be configurable
		err = substitutions.SubstituteFile(pathPrefix+path, file)
		return err
	})

	return err
}

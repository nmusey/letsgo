package generate

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/nmusey/letsgo/lib/files"
	"github.com/nmusey/letsgo/lib/templates"
)

func Run(dir string, kindNames []string, rawName string) error {
	moduleName, err := ModuleName(dir)
	if err != nil {
		return err
	}

	name := NewName(rawName)
	outDir := filepath.Join(dir, "lib", name.Package)

	substitutions := map[string]string{
		"Struct":     name.Struct,
		"Package":    name.Package,
		"ModuleName": moduleName,
	}

	generatorsFS := templates.GeneratorsFilesystem()

	var failures []string
	for _, kindName := range kindNames {
		kind, ok := KindByName(kindName)
		if !ok {
			return fmt.Errorf("generate: unknown kind %q", kindName)
		}

		templatePath := path.Join("generators", kind.TemplateFile)
		outPath := filepath.Join(outDir, kind.OutputFile)

		if err := files.RenderTemplate(generatorsFS, templatePath, outPath, substitutions); err != nil {
			failures = append(failures, fmt.Sprintf("%s: %v", outPath, err))
		}
	}

	if len(failures) > 0 {
		return fmt.Errorf("generate: %s", strings.Join(failures, "; "))
	}

	return nil
}

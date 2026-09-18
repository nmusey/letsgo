package templates_test

import (
	"testing"

	"github.com/nmusey/letsgo/lib/templates"
)

func TestGeneratorsFilesystem_ContainsGeneratorTemplates(t *testing.T) {
	fsys := templates.GeneratorsFilesystem()

	for _, path := range []string{
		"generators/model.go.template",
		"generators/service.go.template",
		"generators/appservice.go.template",
		"generators/repository.go.template",
	} {
		contents, err := fsys.ReadFile(path)
		if err != nil {
			t.Errorf("reading %s: %v", path, err)
			continue
		}
		if len(contents) == 0 {
			t.Errorf("%s is empty", path)
		}
	}
}

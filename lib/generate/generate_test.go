package generate_test

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nmusey/letsgo/lib/generate"
)

func TestRun_WritesRequestedKinds(t *testing.T) {
	dir := t.TempDir()
	writeGoMod(t, dir, "github.com/nmusey/myapp")

	err := generate.Run(dir, []string{"model", "repository", "service"}, "BlogPost")
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	domainDir := filepath.Join(dir, "lib", "blogpost")

	modelContents := readFile(t, filepath.Join(domainDir, "model.go"))
	if !strings.Contains(modelContents, "package blogpost") || !strings.Contains(modelContents, "type BlogPost struct") {
		t.Errorf("model.go = %q, missing expected package/struct", modelContents)
	}
	if !strings.Contains(modelContents, "github.com/nmusey/myapp/lib/database") {
		t.Errorf("model.go = %q, missing module import", modelContents)
	}

	repositoryContents := readFile(t, filepath.Join(domainDir, "repository.go"))
	if !strings.Contains(repositoryContents, "database.Repository[BlogPost]") {
		t.Errorf("repository.go = %q, missing generic repository alias", repositoryContents)
	}

	serviceContents := readFile(t, filepath.Join(domainDir, "service.go"))
	if !strings.Contains(serviceContents, "package blogpost") || !strings.Contains(serviceContents, "type Service struct") {
		t.Errorf("service.go = %q, missing expected package/struct", serviceContents)
	}

	for _, file := range []string{"model.go", "repository.go", "service.go"} {
		fset := token.NewFileSet()
		if _, err := parser.ParseFile(fset, filepath.Join(domainDir, file), nil, 0); err != nil {
			t.Errorf("%s is not syntactically valid Go: %v", file, err)
		}
	}
}

func TestRun_CollisionSkipsOnlyThatFileAndReturnsError(t *testing.T) {
	dir := t.TempDir()
	writeGoMod(t, dir, "github.com/nmusey/myapp")

	domainDir := filepath.Join(dir, "lib", "user")
	if err := os.MkdirAll(domainDir, 0777); err != nil {
		t.Fatalf("seeding domain dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(domainDir, "model.go"), []byte("existing"), 0666); err != nil {
		t.Fatalf("seeding existing model.go: %v", err)
	}

	err := generate.Run(dir, []string{"model", "repository"}, "User")
	if err == nil {
		t.Fatal("expected an error from the model.go collision, got nil")
	}

	if got := readFile(t, filepath.Join(domainDir, "model.go")); got != "existing" {
		t.Errorf("existing model.go was overwritten: got %q", got)
	}

	if _, err := os.Stat(filepath.Join(domainDir, "repository.go")); err != nil {
		t.Errorf("repository.go was not written despite model.go's collision: %v", err)
	}
}

func TestRun_NoGoMod(t *testing.T) {
	dir := t.TempDir()

	if err := generate.Run(dir, []string{"model"}, "User"); err == nil {
		t.Fatal("expected an error when go.mod is missing, got nil")
	}
}

func writeGoMod(t *testing.T, dir, moduleName string) {
	t.Helper()
	content := "module " + moduleName + "\n\ngo 1.26.0\n"
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(content), 0666); err != nil {
		t.Fatalf("seeding go.mod: %v", err)
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	contents, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading %s: %v", path, err)
	}
	return string(contents)
}

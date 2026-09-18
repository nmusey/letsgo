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

	err := generate.Run(dir, []string{"model", "service", "app-service", "repository"}, "BlogPost")
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	modelPath := filepath.Join(dir, "lib", "domain", "blogpost", "model.go")
	modelContents := readFile(t, modelPath)
	if !strings.Contains(modelContents, "package blogpost") || !strings.Contains(modelContents, "type BlogPost struct") {
		t.Errorf("model.go = %q, missing expected package/struct", modelContents)
	}
	if strings.Contains(modelContents, "bun") {
		t.Errorf("model.go = %q, should not reference bun", modelContents)
	}

	domainServicePath := filepath.Join(dir, "lib", "domain", "blogpost", "service.go")
	domainServiceContents := readFile(t, domainServicePath)
	if !strings.Contains(domainServiceContents, "package blogpost") || !strings.Contains(domainServiceContents, "type Service struct") {
		t.Errorf("domain service.go = %q, missing expected package/struct", domainServiceContents)
	}
	if strings.Contains(domainServiceContents, "Repository") {
		t.Errorf("domain service.go = %q, should not reference a repository", domainServiceContents)
	}

	appServicePath := filepath.Join(dir, "lib", "application", "blogpost", "blogpostservice.go")
	appServiceContents := readFile(t, appServicePath)
	if !strings.Contains(appServiceContents, "package blogpost") || !strings.Contains(appServiceContents, "type Repository interface") {
		t.Errorf("appservice.go = %q, missing expected package/interface", appServiceContents)
	}
	if !strings.Contains(appServiceContents, "github.com/nmusey/myapp/lib/domain/blogpost") {
		t.Errorf("appservice.go = %q, missing domain import", appServiceContents)
	}

	repositoryPath := filepath.Join(dir, "lib", "infrastructure", "database", "blogpost", "repository.go")
	repositoryContents := readFile(t, repositoryPath)
	if !strings.Contains(repositoryContents, "database.Repository[SavedBlogPost]") {
		t.Errorf("repository.go = %q, missing generic repository engine", repositoryContents)
	}
	if !strings.Contains(repositoryContents, "github.com/nmusey/myapp/lib/domain/blogpost") {
		t.Errorf("repository.go = %q, missing domain import", repositoryContents)
	}

	for _, path := range []string{modelPath, domainServicePath, appServicePath, repositoryPath} {
		fset := token.NewFileSet()
		if _, err := parser.ParseFile(fset, path, nil, 0); err != nil {
			t.Errorf("%s is not syntactically valid Go: %v", path, err)
		}
	}
}

func TestRun_CollisionSkipsOnlyThatFileAndReturnsError(t *testing.T) {
	dir := t.TempDir()
	writeGoMod(t, dir, "github.com/nmusey/myapp")

	domainDir := filepath.Join(dir, "lib", "domain", "user")
	if err := os.MkdirAll(domainDir, 0777); err != nil {
		t.Fatalf("seeding domain dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(domainDir, "model.go"), []byte("existing"), 0666); err != nil {
		t.Fatalf("seeding existing model.go: %v", err)
	}

	err := generate.Run(dir, []string{"model", "service"}, "User")
	if err == nil {
		t.Fatal("expected an error from the model.go collision, got nil")
	}

	if got := readFile(t, filepath.Join(domainDir, "model.go")); got != "existing" {
		t.Errorf("existing model.go was overwritten: got %q", got)
	}

	if _, err := os.Stat(filepath.Join(domainDir, "service.go")); err != nil {
		t.Errorf("service.go was not written despite model.go's collision: %v", err)
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

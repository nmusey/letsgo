package cli_test

import (
	"context"
	"go/parser"
	"go/token"
	"io"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/nmusey/letsgo/lib/cli"
)

func run(args ...string) error {
	cmd := cli.MakeCommands()
	cmd.Writer = io.Discard
	cmd.ErrWriter = io.Discard
	return cmd.Run(context.Background(), append([]string{"letsgo"}, args...))
}

func writeGoMod(t *testing.T, dir, moduleName string) {
	t.Helper()
	content := "module " + moduleName + "\n\ngo 1.26.0\n"
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(content), 0666); err != nil {
		t.Fatalf("seeding go.mod: %v", err)
	}
}

func newModule(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	writeGoMod(t, dir, "github.com/nmusey/myapp")
	t.Chdir(dir)
	return dir
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	contents, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading %s: %v", path, err)
	}
	return string(contents)
}

func TestMakeProject_GeneratesProjectInNamedDirectory(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)

	if err := run("make", "project", "myapp", "github.com/nmusey/myapp"); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	for _, file := range []string{
		"go.mod",
		"Makefile",
		"docker-compose.yml",
		filepath.Join("cmd", "http", "main.go"),
		filepath.Join("lib", "infrastructure", "database", "connection.go"),
		filepath.Join("lib", "infrastructure", "database", "repository.go"),
	} {
		if _, err := os.Stat(filepath.Join(dir, "myapp", file)); err != nil {
			t.Errorf("expected %s to be created: %v", file, err)
		}
	}
	if got := readFile(t, filepath.Join(dir, "myapp", "go.mod")); !strings.Contains(got, "module github.com/nmusey/myapp") {
		t.Errorf("go.mod does not declare the requested module:\n%s", got)
	}
}

func TestMakeProject_DefaultsDatabasePort(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)

	if err := run("make", "project", "myapp", "github.com/nmusey/myapp"); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	if got := readFile(t, filepath.Join(dir, "myapp", "docker-compose.yml")); !strings.Contains(got, "${DB_PORT}:5432") {
		t.Errorf("docker-compose.yml should map to the default port 5432:\n%s", got)
	}
}

func TestMakeProject_UsesGivenDatabasePort(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)

	if err := run("make", "project", "myapp", "github.com/nmusey/myapp", "6543"); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	if got := readFile(t, filepath.Join(dir, "myapp", "docker-compose.yml")); !strings.Contains(got, "${DB_PORT}:6543") {
		t.Errorf("docker-compose.yml should map to port 6543:\n%s", got)
	}
}

func TestMakeProject_MissingArguments_ReturnsError(t *testing.T) {
	t.Chdir(t.TempDir())

	for _, args := range [][]string{
		{"make", "project"},
		{"make", "project", "myapp"},
	} {
		if err := run(args...); err == nil {
			t.Errorf("run(%v) = nil, want an error for missing required arguments", args)
		}
	}
}

func TestMakeModel_GeneratesModelRepositoryAndService(t *testing.T) {
	dir := newModule(t)

	if err := run("make", "model", "User"); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	for _, file := range []string{"model.go", "repository.go", "service.go"} {
		if _, err := os.Stat(filepath.Join(dir, "lib", "domain", "user", file)); err != nil {
			t.Errorf("expected %s to be created: %v", file, err)
		}
	}
}

func TestMakeModel_NormalisesNames(t *testing.T) {
	dir := newModule(t)

	if err := run("make", "model", "user profile"); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	model := readFile(t, filepath.Join(dir, "lib", "domain", "userprofile", "model.go"))
	if !strings.Contains(model, "package userprofile") || !strings.Contains(model, "type UserProfile struct") {
		t.Errorf("model.go does not use the normalised package and struct name:\n%s", model)
	}
}

func TestMakeModel_AliasMake(t *testing.T) {
	dir := newModule(t)

	if err := run("m", "model", "Post"); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	if _, err := os.Stat(filepath.Join(dir, "lib", "domain", "post", "model.go")); err != nil {
		t.Errorf("expected model.go to be created through the alias: %v", err)
	}
}

func TestMakeModel_NoName_ReturnsError(t *testing.T) {
	newModule(t)

	if err := run("make", "model"); err == nil {
		t.Fatal("expected an error when no model name is given, got nil")
	}
}

func TestMakeModel_OutsideModule_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)

	if err := run("make", "model", "User"); err == nil {
		t.Fatal("expected an error outside of a go module, got nil")
	}
	if _, err := os.Stat(filepath.Join(dir, "lib")); err == nil {
		t.Error("nothing should be generated outside of a go module")
	}
}

func TestMakeAppService_GeneratesFile(t *testing.T) {
	dir := newModule(t)

	if err := run("make", "app-service", "Order"); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	if _, err := os.Stat(filepath.Join(dir, "lib", "application", "order", "appservice.go")); err != nil {
		t.Errorf("expected appservice.go to be created: %v", err)
	}
}

func TestMakeAppService_NoName_ReturnsError(t *testing.T) {
	newModule(t)

	if err := run("make", "app-service"); err == nil {
		t.Fatal("expected an error when no service name is given, got nil")
	}
}

func TestMakeAppService_OutsideModule_ReturnsError(t *testing.T) {
	t.Chdir(t.TempDir())

	if err := run("make", "app-service", "Order"); err == nil {
		t.Fatal("expected an error outside of a go module, got nil")
	}
}

// A package that imports itself is an import cycle, so the generated app service must not live in the package of
// the model it wraps.
func TestMakeAppService_DoesNotImportItsOwnPackage(t *testing.T) {
	dir := newModule(t)

	if err := run("make", "app-service", "Order"); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	generated := filepath.Join(dir, "lib", "application", "order", "appservice.go")
	file, err := parser.ParseFile(token.NewFileSet(), generated, nil, parser.ImportsOnly)
	if err != nil {
		t.Fatalf("parsing %s: %v", generated, err)
	}

	ownImportPath := path.Join("github.com/nmusey/myapp", "lib", "application", "order")
	for _, spec := range file.Imports {
		imported, err := strconv.Unquote(spec.Path.Value)
		if err != nil {
			t.Fatal(err)
		}
		if imported == ownImportPath {
			t.Errorf("appservice.go imports its own package %q", imported)
		}
	}
}

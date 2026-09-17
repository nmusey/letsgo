package generate_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nmusey/letsgo/lib/generate"
)

func TestModuleName(t *testing.T) {
	dir := t.TempDir()
	content := "module github.com/nmusey/myapp\n\ngo 1.26.0\n\nrequire (\n\tgithub.com/urfave/cli/v3 v3.10.0\n)\n"
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(content), 0666); err != nil {
		t.Fatalf("seeding go.mod: %v", err)
	}

	got, err := generate.ModuleName(dir)
	if err != nil {
		t.Fatalf("ModuleName returned error: %v", err)
	}

	want := "github.com/nmusey/myapp"
	if got != want {
		t.Errorf("ModuleName() = %q, want %q", got, want)
	}
}

func TestModuleName_MissingGoMod(t *testing.T) {
	dir := t.TempDir()

	if _, err := generate.ModuleName(dir); err == nil {
		t.Fatal("expected an error when go.mod is missing, got nil")
	}
}

func TestModuleName_NoModuleDirective(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("go 1.26.0\n"), 0666); err != nil {
		t.Fatalf("seeding go.mod: %v", err)
	}

	if _, err := generate.ModuleName(dir); err == nil {
		t.Fatal("expected an error when go.mod has no module directive, got nil")
	}
}

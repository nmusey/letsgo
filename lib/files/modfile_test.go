package files_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nmusey/letsgo/lib/files"
)

func writeGoMod(t *testing.T, dir, contents string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(contents), 0644); err != nil {
		t.Fatalf("seeding go.mod: %v", err)
	}
}

func TestModuleName_GoModInCurrentDirectory(t *testing.T) {
	dir := t.TempDir()
	writeGoMod(t, dir, "module example.com/app\n\ngo 1.26.0\n")
	t.Chdir(dir)

	got, err := files.ModuleName()
	if err != nil {
		t.Fatalf("ModuleName returned error: %v", err)
	}
	if got != "example.com/app" {
		t.Errorf("ModuleName() = %q, want %q", got, "example.com/app")
	}
}

func TestModuleName_GoModInParentDirectory(t *testing.T) {
	dir := t.TempDir()
	writeGoMod(t, dir, "module example.com/app\n")
	nested := filepath.Join(dir, "lib", "domain")
	if err := os.MkdirAll(nested, 0755); err != nil {
		t.Fatal(err)
	}
	t.Chdir(nested)

	got, err := files.ModuleName()
	if err != nil {
		t.Fatalf("ModuleName returned error: %v", err)
	}
	if got != "example.com/app" {
		t.Errorf("ModuleName() = %q, want %q", got, "example.com/app")
	}
}

func TestModuleName_NearestGoModWins(t *testing.T) {
	outer := t.TempDir()
	writeGoMod(t, outer, "module example.com/outer\n")
	inner := filepath.Join(outer, "inner")
	if err := os.Mkdir(inner, 0755); err != nil {
		t.Fatal(err)
	}
	writeGoMod(t, inner, "module example.com/inner\n")
	t.Chdir(inner)

	got, err := files.ModuleName()
	if err != nil {
		t.Fatalf("ModuleName returned error: %v", err)
	}
	if got != "example.com/inner" {
		t.Errorf("ModuleName() = %q, want %q", got, "example.com/inner")
	}
}

func TestModuleName_DirectoryNamedGoModIsSkipped(t *testing.T) {
	outer := t.TempDir()
	writeGoMod(t, outer, "module example.com/outer\n")
	inner := filepath.Join(outer, "inner")
	if err := os.MkdirAll(filepath.Join(inner, "go.mod"), 0755); err != nil {
		t.Fatal(err)
	}
	t.Chdir(inner)

	got, err := files.ModuleName()
	if err != nil {
		t.Fatalf("ModuleName returned error: %v", err)
	}
	if got != "example.com/outer" {
		t.Errorf("ModuleName() = %q, want %q", got, "example.com/outer")
	}
}

func TestModuleName_NoGoMod_ReturnsError(t *testing.T) {
	t.Chdir(t.TempDir())

	_, err := files.ModuleName()
	if err == nil {
		t.Fatal("expected an error when no go.mod exists, got nil")
	}
	if !strings.Contains(err.Error(), "no go.mod found") {
		t.Errorf("error = %q, want it to mention that no go.mod was found", err.Error())
	}
}

func TestModuleName_NoModuleDeclaration_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	writeGoMod(t, dir, "go 1.26.0\n")
	t.Chdir(dir)

	_, err := files.ModuleName()
	if err == nil {
		t.Fatal("expected an error when go.mod has no module line, got nil")
	}
	if !strings.Contains(err.Error(), "no module declaration") {
		t.Errorf("error = %q, want it to mention the missing module declaration", err.Error())
	}
}

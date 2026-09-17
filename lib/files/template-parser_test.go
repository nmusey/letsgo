package files_test

import (
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"

	"github.com/nmusey/letsgo/lib/files"
)

func TestRenderTemplate_WritesSubstitutedFile(t *testing.T) {
	fsys := fstest.MapFS{
		"generators/model.go.template": &fstest.MapFile{
			Data: []byte("package {{.Package}}\n\ntype {{.Struct}} struct{}\n"),
		},
	}

	dir := t.TempDir()
	outPath := filepath.Join(dir, "lib", "user", "model.go")

	err := files.RenderTemplate(fsys, "generators/model.go.template", outPath, map[string]string{
		"Package": "user",
		"Struct":  "User",
	})
	if err != nil {
		t.Fatalf("RenderTemplate returned error: %v", err)
	}

	got, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("reading generated file: %v", err)
	}

	want := "package user\n\ntype User struct{}\n"
	if string(got) != want {
		t.Errorf("generated file = %q, want %q", got, want)
	}
}

func TestRenderTemplate_RefusesToOverwriteExistingFile(t *testing.T) {
	fsys := fstest.MapFS{
		"generators/model.go.template": &fstest.MapFile{
			Data: []byte("package {{.Package}}\n"),
		},
	}

	dir := t.TempDir()
	outPath := filepath.Join(dir, "model.go")

	if err := os.WriteFile(outPath, []byte("existing content"), 0666); err != nil {
		t.Fatalf("seeding existing file: %v", err)
	}

	err := files.RenderTemplate(fsys, "generators/model.go.template", outPath, map[string]string{
		"Package": "user",
	})
	if err == nil {
		t.Fatal("expected an error when outPath already exists, got nil")
	}

	got, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("reading file after failed RenderTemplate: %v", err)
	}
	if string(got) != "existing content" {
		t.Errorf("existing file was modified: got %q", got)
	}
}

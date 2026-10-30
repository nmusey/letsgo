package files_test

import (
	"bytes"
	"embed"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nmusey/letsgo/lib/files"
)

//go:embed all:testdata/fixture
var fixtureFS embed.FS

//go:embed all:testdata/broken
var brokenFS embed.FS

func fixtureParser() files.TemplateParser {
	return files.TemplateParser{
		Filesystem:     fixtureFS,
		FilesystemRoot: "testdata/fixture",
		Substitutions:  map[string]string{"Name": "World"},
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

func TestParseFiles_SubstitutesTemplatesAndStripsSuffix(t *testing.T) {
	out := filepath.Join(t.TempDir(), "out")

	if err := fixtureParser().ParseFiles(out); err != nil {
		t.Fatalf("ParseFiles returned error: %v", err)
	}

	if got, want := readFile(t, filepath.Join(out, "greeting.txt")), "Hello World!\n"; got != want {
		t.Errorf("greeting.txt = %q, want %q", got, want)
	}
	if _, err := os.Stat(filepath.Join(out, "greeting.txt.template")); err == nil {
		t.Error("the .template suffix should have been stripped from the output name")
	}
}

func TestParseFiles_CreatesNestedDirectories(t *testing.T) {
	out := filepath.Join(t.TempDir(), "out")

	if err := fixtureParser().ParseFiles(out); err != nil {
		t.Fatalf("ParseFiles returned error: %v", err)
	}

	info, err := os.Stat(filepath.Join(out, "nested"))
	if err != nil || !info.IsDir() {
		t.Fatalf("expected nested directory to be created, stat err = %v", err)
	}
	if got, want := readFile(t, filepath.Join(out, "nested", "plain.txt")), "no World substitution here\n"; got != want {
		t.Errorf("nested/plain.txt = %q, want %q", got, want)
	}
}

func TestParseFiles_IncludesHiddenFiles(t *testing.T) {
	out := filepath.Join(t.TempDir(), "out")

	if err := fixtureParser().ParseFiles(out); err != nil {
		t.Fatalf("ParseFiles returned error: %v", err)
	}

	if got, want := readFile(t, filepath.Join(out, ".hidden")), "hidden World\n"; got != want {
		t.Errorf(".hidden = %q, want %q", got, want)
	}
}

func TestParseFiles_CreatesMissingParentDirectories(t *testing.T) {
	out := filepath.Join(t.TempDir(), "a", "b", "c")

	if err := fixtureParser().ParseFiles(out); err != nil {
		t.Fatalf("ParseFiles returned error: %v", err)
	}

	if _, err := os.Stat(filepath.Join(out, "greeting.txt")); err != nil {
		t.Errorf("expected greeting.txt in nested output directory: %v", err)
	}
}

func TestParseFiles_OverwritesExistingFiles(t *testing.T) {
	out := filepath.Join(t.TempDir(), "out")
	if err := os.MkdirAll(out, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(out, "greeting.txt"), []byte("stale"), 0644); err != nil {
		t.Fatal(err)
	}

	if err := fixtureParser().ParseFiles(out); err != nil {
		t.Fatalf("ParseFiles returned error: %v", err)
	}

	if got, want := readFile(t, filepath.Join(out, "greeting.txt")), "Hello World!\n"; got != want {
		t.Errorf("greeting.txt = %q, want %q", got, want)
	}
}

func TestParseFiles_UnwritableOutput_ReturnsError(t *testing.T) {
	blocker := filepath.Join(t.TempDir(), "blocker")
	if err := os.WriteFile(blocker, []byte("a file, not a directory"), 0644); err != nil {
		t.Fatal(err)
	}

	if err := fixtureParser().ParseFiles(filepath.Join(blocker, "out")); err == nil {
		t.Fatal("expected an error when the output directory cannot be created, got nil")
	}
}

func TestParseFiles_InvalidTemplate_ReturnsError(t *testing.T) {
	parser := files.TemplateParser{
		Filesystem:     brokenFS,
		FilesystemRoot: "testdata/broken",
		Substitutions:  map[string]string{"Name": "World"},
	}

	if err := parser.ParseFiles(filepath.Join(t.TempDir(), "out")); err == nil {
		t.Fatal("expected an error for a template that does not parse, got nil")
	}
}

func TestSubstituteFile_WritesSubstitutedContents(t *testing.T) {
	var buf bytes.Buffer

	if err := fixtureParser().SubstituteFile("testdata/fixture/greeting.txt.template", &buf); err != nil {
		t.Fatalf("SubstituteFile returned error: %v", err)
	}
	if got, want := buf.String(), "Hello World!\n"; got != want {
		t.Errorf("output = %q, want %q", got, want)
	}
}

func TestSubstituteFile_MissingFile_ReturnsError(t *testing.T) {
	var buf bytes.Buffer

	if err := fixtureParser().SubstituteFile("testdata/fixture/does-not-exist", &buf); err == nil {
		t.Fatal("expected an error for a file that does not exist, got nil")
	}
}

func TestSubstituteFile_InvalidTemplate_ReturnsError(t *testing.T) {
	parser := files.TemplateParser{Filesystem: brokenFS, FilesystemRoot: "testdata/broken"}
	var buf bytes.Buffer

	err := parser.SubstituteFile("testdata/broken/bad.template", &buf)
	if err == nil {
		t.Fatal("expected a parse error, got nil")
	}
	if !strings.Contains(err.Error(), "bad.template") {
		t.Errorf("error = %q, want it to name the offending template", err.Error())
	}
}

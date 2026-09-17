package cli_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nmusey/letsgo/lib/cli"
)

func TestMakeModel_WithCompanionFlags_GeneratesAllThree(t *testing.T) {
	dir := t.TempDir()
	writeGoMod(t, dir, "github.com/nmusey/myapp")
	t.Chdir(dir)

	cmd := cli.MakeCommands()
	if err := cmd.Run(context.Background(), []string{"letsgo", "make", "model", "User", "-r", "-s"}); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	for _, file := range []string{"model.go", "repository.go", "service.go"} {
		path := filepath.Join(dir, "lib", "user", file)
		if _, err := os.Stat(path); err != nil {
			t.Errorf("expected %s to be created: %v", path, err)
		}
	}
}

func TestMakeRepository_WithServiceFlag_DoesNotGenerateModel(t *testing.T) {
	dir := t.TempDir()
	writeGoMod(t, dir, "github.com/nmusey/myapp")
	t.Chdir(dir)

	cmd := cli.MakeCommands()
	if err := cmd.Run(context.Background(), []string{"letsgo", "make", "repository", "Post", "-s"}); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	domainDir := filepath.Join(dir, "lib", "post")
	if _, err := os.Stat(filepath.Join(domainDir, "repository.go")); err != nil {
		t.Errorf("expected repository.go to be created: %v", err)
	}
	if _, err := os.Stat(filepath.Join(domainDir, "service.go")); err != nil {
		t.Errorf("expected service.go to be created: %v", err)
	}
	if _, err := os.Stat(filepath.Join(domainDir, "model.go")); err == nil {
		t.Error("model.go should not have been created")
	}
}

func TestMakeService_NoName_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	writeGoMod(t, dir, "github.com/nmusey/myapp")
	t.Chdir(dir)

	cmd := cli.MakeCommands()
	err := cmd.Run(context.Background(), []string{"letsgo", "make", "service"})
	if err == nil {
		t.Fatal("expected an error when no domain name is given, got nil")
	}
	if !strings.Contains(err.Error(), "domain name is required") {
		t.Errorf("error = %q, want it to mention a required domain name", err.Error())
	}
}

func writeGoMod(t *testing.T, dir, moduleName string) {
	t.Helper()
	content := "module " + moduleName + "\n\ngo 1.26.0\n"
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(content), 0666); err != nil {
		t.Fatalf("seeding go.mod: %v", err)
	}
}

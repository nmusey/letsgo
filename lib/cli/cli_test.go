package cli_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nmusey/letsgo/lib/cli"
)

func TestMakeModel_WithCompanionFlags_GeneratesAllFour(t *testing.T) {
	dir := t.TempDir()
	writeGoMod(t, dir, "github.com/nmusey/myapp")
	t.Chdir(dir)

	cmd := cli.MakeCommands()
	if err := cmd.Run(context.Background(), []string{"letsgo", "make", "model", "User", "-s", "-a", "-r"}); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	for _, path := range []string{
		filepath.Join(dir, "lib", "domain", "user", "model.go"),
		filepath.Join(dir, "lib", "domain", "user", "service.go"),
		filepath.Join(dir, "lib", "application", "user", "userservice.go"),
		filepath.Join(dir, "lib", "infrastructure", "database", "user", "repository.go"),
	} {
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

	if _, err := os.Stat(filepath.Join(dir, "lib", "infrastructure", "database", "post", "repository.go")); err != nil {
		t.Errorf("expected repository.go to be created: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "lib", "domain", "post", "service.go")); err != nil {
		t.Errorf("expected service.go to be created: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "lib", "domain", "post", "model.go")); err == nil {
		t.Error("model.go should not have been created")
	}
	if _, err := os.Stat(filepath.Join(dir, "lib", "application", "post", "postservice.go")); err == nil {
		t.Error("postservice.go should not have been created")
	}
}

func TestMakeAppService_WithModelFlag_GeneratesBoth(t *testing.T) {
	dir := t.TempDir()
	writeGoMod(t, dir, "github.com/nmusey/myapp")
	t.Chdir(dir)

	cmd := cli.MakeCommands()
	if err := cmd.Run(context.Background(), []string{"letsgo", "make", "app-service", "Order", "-m"}); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	if _, err := os.Stat(filepath.Join(dir, "lib", "application", "order", "orderservice.go")); err != nil {
		t.Errorf("expected orderservice.go to be created: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "lib", "domain", "order", "model.go")); err != nil {
		t.Errorf("expected model.go to be created: %v", err)
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

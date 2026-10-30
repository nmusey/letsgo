package templates

import (
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

const databaseModule = "github.com/nmusey/letsgo/lib/templates/project"

func writeGoMod(t *testing.T, dir, moduleName string) {
	t.Helper()
	content := "module " + moduleName + "\n\ngo 1.26.0\n"
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(content), 0644); err != nil {
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

func TestNewProjectTemplate_SetsSubstitutions(t *testing.T) {
	parser := NewProjectTemplate("myapp", "example.com/myapp", "6543")

	want := map[string]string{"AppName": "myapp", "Module": "example.com/myapp", "DbPort": "6543"}
	for key, value := range want {
		if parser.Substitutions[key] != value {
			t.Errorf("Substitutions[%q] = %q, want %q", key, parser.Substitutions[key], value)
		}
	}
	if parser.FilesystemRoot != "project" {
		t.Errorf("FilesystemRoot = %q, want %q", parser.FilesystemRoot, "project")
	}
}

func TestNewProjectTemplate_GeneratesProject(t *testing.T) {
	out := filepath.Join(t.TempDir(), "myapp")

	if err := NewProjectTemplate("myapp", "example.com/myapp", "6543").ParseFiles(out); err != nil {
		t.Fatalf("ParseFiles returned error: %v", err)
	}

	for _, path := range []string{
		".air.toml",
		".env.example",
		".gitignore",
		"Dockerfile",
		"Makefile",
		"docker-compose.yml",
		"go.mod",
		filepath.Join("cmd", "http", "main.go"),
		filepath.Join("lib", "infrastructure", "database", "connection.go"),
		filepath.Join("lib", "infrastructure", "database", "repository.go"),
	} {
		if _, err := os.Stat(filepath.Join(out, path)); err != nil {
			t.Errorf("expected %s to be generated: %v", path, err)
		}
	}

	assertContains(t, readFile(t, filepath.Join(out, "go.mod")), "module example.com/myapp")
	assertContains(t, readFile(t, filepath.Join(out, "Makefile")), "myapp:dev")
	assertContains(t, readFile(t, filepath.Join(out, "docker-compose.yml")), "${DB_PORT}:6543")
	assertContains(t, readFile(t, filepath.Join(out, "cmd", "http", "main.go")), `"myapp/lib/infrastructure/database"`)
}

func TestNewProjectTemplate_LeavesNoUnrenderedActions(t *testing.T) {
	out := filepath.Join(t.TempDir(), "myapp")

	if err := NewProjectTemplate("myapp", "example.com/myapp", "5432").ParseFiles(out); err != nil {
		t.Fatalf("ParseFiles returned error: %v", err)
	}

	err := filepath.WalkDir(out, func(path string, entry os.DirEntry, err error) error {
		if err != nil || entry.IsDir() {
			return err
		}
		if strings.Contains(readFile(t, path), "{{") {
			t.Errorf("%s still contains an unrendered template action", path)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewProjectTemplate_GeneratedGoFilesParse(t *testing.T) {
	out := filepath.Join(t.TempDir(), "myapp")

	if err := NewProjectTemplate("myapp", "example.com/myapp", "5432").ParseFiles(out); err != nil {
		t.Fatalf("ParseFiles returned error: %v", err)
	}

	assertGoFilesParse(t, out)
}

func TestNewModelTemplate_SetsSubstitutions(t *testing.T) {
	dir := t.TempDir()
	writeGoMod(t, dir, "example.com/app")
	t.Chdir(dir)

	parser, err := NewModelTemplate("user profile")
	if err != nil {
		t.Fatalf("NewModelTemplate returned error: %v", err)
	}

	want := map[string]string{"Package": "userprofile", "Struct": "UserProfile", "Module": "example.com/app"}
	for key, value := range want {
		if parser.Substitutions[key] != value {
			t.Errorf("Substitutions[%q] = %q, want %q", key, parser.Substitutions[key], value)
		}
	}
	if parser.FilesystemRoot != "model" {
		t.Errorf("FilesystemRoot = %q, want %q", parser.FilesystemRoot, "model")
	}
}

func TestNewModelTemplate_NoGoMod_ReturnsError(t *testing.T) {
	t.Chdir(t.TempDir())

	if _, err := NewModelTemplate("user"); err == nil {
		t.Fatal("expected an error outside of a go module, got nil")
	}
}

func TestNewModelTemplate_GeneratesFiles(t *testing.T) {
	dir := t.TempDir()
	writeGoMod(t, dir, "example.com/app")
	t.Chdir(dir)

	parser, err := NewModelTemplate("user profile")
	if err != nil {
		t.Fatal(err)
	}
	out := filepath.Join(dir, "lib", "domain", "userprofile")
	if err := parser.ParseFiles(out); err != nil {
		t.Fatalf("ParseFiles returned error: %v", err)
	}

	model := readFile(t, filepath.Join(out, "model.go"))
	assertContains(t, model, "package userprofile")
	assertContains(t, model, "type UserProfile struct")

	repository := readFile(t, filepath.Join(out, "repository.go"))
	assertContains(t, repository, "package userprofile")
	assertContains(t, repository, `"example.com/app/lib/infrastructure/database"`)
	assertContains(t, repository, "type SavedUserProfile struct")
	assertContains(t, repository, `bun:"table:userprofile"`)

	assertContains(t, readFile(t, filepath.Join(out, "service.go")), "package userprofile")
}

func TestNewModelTemplate_GeneratedGoFilesParse(t *testing.T) {
	dir := t.TempDir()
	writeGoMod(t, dir, "example.com/app")
	t.Chdir(dir)

	parser, err := NewModelTemplate("user profile")
	if err != nil {
		t.Fatal(err)
	}
	out := filepath.Join(dir, "out")
	if err := parser.ParseFiles(out); err != nil {
		t.Fatalf("ParseFiles returned error: %v", err)
	}

	assertGoFilesParse(t, out)
}

func TestNewModelTemplate_RepositoryDoesNotQualifyOwnPackage(t *testing.T) {
	dir := t.TempDir()
	writeGoMod(t, dir, "example.com/app")
	t.Chdir(dir)

	parser, err := NewModelTemplate("user")
	if err != nil {
		t.Fatal(err)
	}
	out := filepath.Join(dir, "out")
	if err := parser.ParseFiles(out); err != nil {
		t.Fatalf("ParseFiles returned error: %v", err)
	}

	if repository := readFile(t, filepath.Join(out, "repository.go")); strings.Contains(repository, "user.User") {
		t.Error("repository.go refers to its own package by name (user.User), which does not compile")
	}
}

func TestNewAppServiceTemplate_SetsSubstitutions(t *testing.T) {
	dir := t.TempDir()
	writeGoMod(t, dir, "example.com/app")
	t.Chdir(dir)

	parser, err := NewAppServiceTemplate("order-item")
	if err != nil {
		t.Fatalf("NewAppServiceTemplate returned error: %v", err)
	}

	want := map[string]string{"Package": "orderitem", "Struct": "OrderItem", "Module": "example.com/app"}
	for key, value := range want {
		if parser.Substitutions[key] != value {
			t.Errorf("Substitutions[%q] = %q, want %q", key, parser.Substitutions[key], value)
		}
	}
	if parser.FilesystemRoot != "appservice" {
		t.Errorf("FilesystemRoot = %q, want %q", parser.FilesystemRoot, "appservice")
	}
}

func TestNewAppServiceTemplate_NoGoMod_ReturnsError(t *testing.T) {
	t.Chdir(t.TempDir())

	if _, err := NewAppServiceTemplate("order"); err == nil {
		t.Fatal("expected an error outside of a go module, got nil")
	}
}

func TestNewAppServiceTemplate_GeneratesFile(t *testing.T) {
	dir := t.TempDir()
	writeGoMod(t, dir, "example.com/app")
	t.Chdir(dir)

	parser, err := NewAppServiceTemplate("order")
	if err != nil {
		t.Fatal(err)
	}
	out := filepath.Join(dir, "out")
	if err := parser.ParseFiles(out); err != nil {
		t.Fatalf("ParseFiles returned error: %v", err)
	}

	contents := readFile(t, filepath.Join(out, "appservice.go"))
	assertContains(t, contents, "package order")
	assertContains(t, contents, `"example.com/app/lib/domain/order"`)
	assertContains(t, contents, "Get(ctx context.Context, id int64) (*order.Order, error)")
	assertGoFilesParse(t, out)
}

// The generated repository is compiled and exercised against an in-memory sqlite database. The output has to live
// inside this module so that it can import the project's database package, which stands in for the one that
// `letsgo make project` generates.
func TestModelTemplate_RepositoryRoundTripsThroughSQLite(t *testing.T) {
	goBinary, err := exec.LookPath("go")
	if err != nil {
		t.Skip("go toolchain not available")
	}

	out := generateModelInModule(t, "user profile")
	// The service template is covered separately, keep it from breaking the build of this package.
	if err := os.Remove(filepath.Join(out, "service.go")); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(out, "roundtrip_test.go"), []byte(roundTripTest), 0644); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(goBinary, "test", "./"+filepath.Base(out))
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("generated package failed its tests: %v\n%s", err, output)
	}
}

func TestModelTemplate_GeneratedPackageCompiles(t *testing.T) {
	goBinary, err := exec.LookPath("go")
	if err != nil {
		t.Skip("go toolchain not available")
	}

	out := generateModelInModule(t, "user profile")

	cmd := exec.Command(goBinary, "vet", "./"+filepath.Base(out))
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("generated model package does not compile: %v\n%s", err, output)
	}
}

func generateModelInModule(t *testing.T, name string) string {
	t.Helper()

	packageDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	// Directories starting with an underscore are skipped by `./...` so this never leaks into other test runs.
	out, err := os.MkdirTemp(packageDir, "_generated")
	if err != nil {
		t.Fatalf("creating output directory: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(out) })

	// The templates read the module name from the nearest go.mod, so render them from a directory that fakes one.
	dir := t.TempDir()
	writeGoMod(t, dir, databaseModule)
	t.Chdir(dir)
	parser, err := NewModelTemplate(name)
	if err != nil {
		t.Fatal(err)
	}
	t.Chdir(packageDir)

	if err := parser.ParseFiles(out); err != nil {
		t.Fatalf("ParseFiles returned error: %v", err)
	}
	return out
}

func assertContains(t *testing.T, contents, want string) {
	t.Helper()
	if !strings.Contains(contents, want) {
		t.Errorf("expected output to contain %q, got:\n%s", want, contents)
	}
}

func assertGoFilesParse(t *testing.T, root string) {
	t.Helper()

	fileSet := token.NewFileSet()
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil || entry.IsDir() || filepath.Ext(path) != ".go" {
			return err
		}
		if _, err := parser.ParseFile(fileSet, path, nil, parser.AllErrors); err != nil {
			t.Errorf("%s is not valid Go: %v", path, err)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

const roundTripTest = `package userprofile

import (
	"context"
	"testing"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"database/sql"
)

func newRepository(t *testing.T) *Repository {
	t.Helper()

	sqldb, err := sql.Open(sqliteshim.ShimName, "file::memory:")
	if err != nil {
		t.Fatal(err)
	}
	sqldb.SetMaxOpenConns(1)
	db := bun.NewDB(sqldb, sqlitedialect.New())
	t.Cleanup(func() { db.Close() })

	if _, err := db.NewCreateTable().Model((*SavedUserProfile)(nil)).Exec(context.Background()); err != nil {
		t.Fatalf("creating table: %v", err)
	}

	return New(db)
}

func TestRepository_CreateFillsGeneratedFields(t *testing.T) {
	repo := newRepository(t)
	model := &UserProfile{}

	if err := repo.Create(context.Background(), model); err != nil {
		t.Fatalf("Create: %v", err)
	}

	if model.ID == 0 || model.UUID.String() == "00000000-0000-0000-0000-000000000000" {
		t.Errorf("expected ID and UUID to be set, got %d / %s", model.ID, model.UUID)
	}
	if model.CreatedAt.IsZero() || model.UpdatedAt.IsZero() {
		t.Error("expected timestamps to be set")
	}
}

func TestRepository_GetReturnsCreatedModel(t *testing.T) {
	repo := newRepository(t)
	created := &UserProfile{}
	if err := repo.Create(context.Background(), created); err != nil {
		t.Fatalf("Create: %v", err)
	}

	got, err := repo.Get(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.ID != created.ID || got.UUID != created.UUID {
		t.Errorf("Get returned %+v, want %+v", got, created)
	}
}

func TestRepository_GetMissingReturnsError(t *testing.T) {
	repo := newRepository(t)

	if _, err := repo.Get(context.Background(), 42); err == nil {
		t.Fatal("expected an error for a missing row")
	}
}

func TestRepository_ListReturnsAllModels(t *testing.T) {
	repo := newRepository(t)
	for i := 0; i < 3; i++ {
		if err := repo.Create(context.Background(), &UserProfile{}); err != nil {
			t.Fatalf("Create: %v", err)
		}
	}

	got, err := repo.List(context.Background())
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(got) != 3 {
		t.Fatalf("List returned %d models, want 3", len(got))
	}
	for _, model := range got {
		if model.ID == 0 || model.CreatedAt.IsZero() {
			t.Errorf("List returned a model without its stored fields: %+v", model)
		}
	}
}

func TestRepository_UpdateBumpsUpdatedAt(t *testing.T) {
	repo := newRepository(t)
	model := &UserProfile{}
	if err := repo.Create(context.Background(), model); err != nil {
		t.Fatalf("Create: %v", err)
	}
	createdAt, before := model.CreatedAt, model.UpdatedAt

	if err := repo.Update(context.Background(), model); err != nil {
		t.Fatalf("Update: %v", err)
	}

	if model.UpdatedAt.Before(before) {
		t.Errorf("UpdatedAt went backwards: %s -> %s", before, model.UpdatedAt)
	}
	if !model.CreatedAt.Equal(createdAt) {
		t.Errorf("CreatedAt changed: %s -> %s", createdAt, model.CreatedAt)
	}
}

func TestRepository_DeleteRemovesModel(t *testing.T) {
	repo := newRepository(t)
	model := &UserProfile{}
	if err := repo.Create(context.Background(), model); err != nil {
		t.Fatalf("Create: %v", err)
	}

	if err := repo.Delete(context.Background(), model); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if _, err := repo.Get(context.Background(), model.ID); err == nil {
		t.Error("expected Get to fail after Delete")
	}
}
`

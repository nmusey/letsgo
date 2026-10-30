package templates_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/nmusey/letsgo/lib/templates/project/lib/infrastructure/database"
	"github.com/uptrace/bun"
)

type widget struct {
	bun.BaseModel `bun:"table:widgets"`

	ID   int64  `bun:",pk,autoincrement"`
	Name string `bun:",notnull"`
}

func useSQLite(t *testing.T) {
	t.Helper()
	t.Setenv("DB_DRIVER", "sqlite")
	t.Setenv("DB_SQLITE_PATH", filepath.Join(t.TempDir(), "test.db"))
}

func newWidgetRepository(t *testing.T) *database.Repository[widget] {
	t.Helper()
	useSQLite(t)

	db, err := database.NewDB()
	if err != nil {
		t.Fatalf("NewDB: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	if _, err := db.NewCreateTable().Model((*widget)(nil)).Exec(context.Background()); err != nil {
		t.Fatalf("creating table: %v", err)
	}

	return database.NewRepository[widget](db)
}

func TestNewDB_SQLiteDriver_Connects(t *testing.T) {
	useSQLite(t)

	db, err := database.NewDB()
	if err != nil {
		t.Fatalf("NewDB returned error: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		t.Errorf("Ping: %v", err)
	}
}

func TestNewDB_SQLiteWithoutPath_UsesInMemoryDatabase(t *testing.T) {
	t.Setenv("DB_DRIVER", "sqlite")
	t.Setenv("DB_SQLITE_PATH", "")

	db, err := database.NewDB()
	if err != nil {
		t.Fatalf("NewDB returned error: %v", err)
	}
	defer db.Close()
}

func TestNewDB_UnknownDriver_ReturnsError(t *testing.T) {
	t.Setenv("DB_DRIVER", "mysql")

	db, err := database.NewDB()
	if err == nil {
		db.Close()
		t.Fatal("expected an error for an unknown driver, got nil")
	}
}

func TestNewDB_DefaultsToPostgres(t *testing.T) {
	t.Setenv("DB_DRIVER", "")
	t.Setenv("DB_HOST", "127.0.0.1")
	t.Setenv("DB_PORT", "1")
	t.Setenv("DB_USER", "user")
	t.Setenv("DB_PASS", "pass")
	t.Setenv("DB_NAME", "app")

	db, err := database.NewDB()
	if err == nil {
		db.Close()
		t.Fatal("expected a connection error from an unreachable postgres, got nil")
	}
}

func TestRepository_CreateThenGet(t *testing.T) {
	repo := newWidgetRepository(t)
	ctx := context.Background()
	created := &widget{Name: "sprocket"}

	if err := repo.Create(ctx, created); err != nil {
		t.Fatalf("Create: %v", err)
	}
	if created.ID == 0 {
		t.Error("expected Create to populate the primary key")
	}

	got, err := repo.Get(ctx, created.ID)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.Name != "sprocket" {
		t.Errorf("Get returned name %q, want %q", got.Name, "sprocket")
	}
}

func TestRepository_GetMissing_ReturnsError(t *testing.T) {
	repo := newWidgetRepository(t)

	if _, err := repo.Get(context.Background(), 99); err == nil {
		t.Fatal("expected an error for a missing row, got nil")
	}
}

func TestRepository_List(t *testing.T) {
	repo := newWidgetRepository(t)
	ctx := context.Background()

	empty, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("List on empty table: %v", err)
	}
	if len(empty) != 0 {
		t.Errorf("List on empty table returned %d rows", len(empty))
	}

	for _, name := range []string{"a", "b"} {
		if err := repo.Create(ctx, &widget{Name: name}); err != nil {
			t.Fatalf("Create: %v", err)
		}
	}

	got, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(got) != 2 {
		t.Errorf("List returned %d rows, want 2", len(got))
	}
}

func TestRepository_Update(t *testing.T) {
	repo := newWidgetRepository(t)
	ctx := context.Background()
	model := &widget{Name: "before"}
	if err := repo.Create(ctx, model); err != nil {
		t.Fatalf("Create: %v", err)
	}

	model.Name = "after"
	if err := repo.Update(ctx, model); err != nil {
		t.Fatalf("Update: %v", err)
	}

	got, err := repo.Get(ctx, model.ID)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.Name != "after" {
		t.Errorf("name = %q after update, want %q", got.Name, "after")
	}
}

func TestRepository_Delete(t *testing.T) {
	repo := newWidgetRepository(t)
	ctx := context.Background()
	model := &widget{Name: "doomed"}
	if err := repo.Create(ctx, model); err != nil {
		t.Fatalf("Create: %v", err)
	}

	if err := repo.Delete(ctx, model); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if _, err := repo.Get(ctx, model.ID); err == nil {
		t.Error("expected Get to fail after Delete")
	}
}

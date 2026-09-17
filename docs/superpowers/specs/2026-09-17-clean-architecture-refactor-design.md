# Clean architecture refactor for generated projects — Design

## Context

The domain-file generators (`letsgo make model|repository|service`, added in
`v2-domain-generators`) currently generate all three files for a domain into
one flat package, `lib/<name>/`: `model.go` (embeds a shared `database.Model`
base), `repository.go` (a generic `database.Repository[T]` alias), and
`service.go` (a bare constructor). The shared generic pieces
(`Model`, `Repository[T]`, `NewDB`) live in `lib/database/` in the scaffolded
project.

This design restructures generated projects around clean/hexagonal
architecture layers, and splits "service" into two distinct roles that the
current design conflates.

## Goals

- **Domain layer** (`lib/domain/<name>/`) holds the pure business model and
  business rules, with zero knowledge of persistence, bun, or SQL.
  - `model.go`: a plain struct (still carries `ID`, `UUID`, `CreatedAt`,
    `UpdatedAt` as ordinary fields — no bun tags, no embedded base).
  - `service.go`: a **domain service** — pure business rules and invariants
    (e.g. "a password must be 8+ characters"). No repository access.
- **Application layer** (`lib/application/<name>/`) holds orchestration.
  - `<name>service.go` (e.g. `userservice.go`): an **application service** —
    declares the `Repository` interface (the port) it depends on, and a
    `Service` that calls through it, dispatches events, etc. Named
    `<name>service.go` rather than `service.go` so it doesn't collide in
    tabs/tooling with the domain layer's own `service.go`.
- **Infrastructure layer** (`lib/infrastructure/`) holds concrete driven
  adapters.
  - `database/repository.go`, `database/connection.go`: the existing generic
    `Repository[T]` engine and `NewDB()`, moved here unchanged (scaffolded
    once per project, not regenerated per domain).
  - `database/<name>/repository.go`: the concrete adapter for one domain.
    Defines `Saved<Struct>` (a bun-tagged wrapper embedding the domain type,
    carrying the UUID/timestamp `BeforeAppendModelHook` previously on the
    shared `Model` base) and a `Repository` that satisfies the application
    layer's `Repository` interface by translating to/from the domain type at
    each method boundary.
- Four independent `letsgo make` kinds — `model`, `service`, `app-service`,
  `repository` — each generatable alone or with any combination of the
  others via companion flags (`-m`/`-s`/`-a`/`-r`), same mechanism as today
  just extended to a fourth kind.
- Dependency direction: `application` and `infrastructure` import `domain`;
  `domain` imports neither.

## Non-goals

- No event-dispatch mechanism is generated — the application service file
  gives a place to add one, nothing more.
- No change to name normalization, collision handling, or the overall
  `letsgo make <kind> <name> [flags]` command shape.
- No automatic mapping code generation beyond the fixed CRUD set
  (Get/List/Create/Update/Delete) already generated today.

## Generated files

**`lib/domain/user/model.go`**
```go
package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        int64
	UUID      uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}
```

**`lib/domain/user/service.go`**
```go
package user

type Service struct {
}

func NewService() *Service {
	return &Service{}
}
```

**`lib/application/user/userservice.go`**
```go
package user

import (
	"context"

	"{{.ModuleName}}/lib/domain/user"
)

type Repository interface {
	Get(ctx context.Context, id int64) (*user.User, error)
	List(ctx context.Context) ([]user.User, error)
	Create(ctx context.Context, model *user.User) error
	Update(ctx context.Context, model *user.User) error
	Delete(ctx context.Context, model *user.User) error
}

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}
```

**`lib/infrastructure/database/user/repository.go`**
```go
package user

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"

	"{{.ModuleName}}/lib/domain/user"
	"{{.ModuleName}}/lib/infrastructure/database"
)

type SavedUser struct {
	bun.BaseModel `bun:"table:user"`

	user.User

	ID        int64     `bun:",pk,autoincrement"`
	UUID      uuid.UUID `bun:"type:uuid,unique,notnull"`
	CreatedAt time.Time `bun:",notnull"`
	UpdatedAt time.Time `bun:",notnull"`
}

var _ bun.BeforeAppendModelHook = (*SavedUser)(nil)

func (m *SavedUser) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		if m.UUID == uuid.Nil {
			m.UUID = uuid.New()
		}
		now := time.Now()
		if m.CreatedAt.IsZero() {
			m.CreatedAt = now
		}
		m.UpdatedAt = now
	case *bun.UpdateQuery:
		m.UpdatedAt = time.Now()
	}
	return nil
}

type Repository struct {
	engine *database.Repository[SavedUser]
}

func NewRepository(db *bun.DB) *Repository {
	return &Repository{engine: database.NewRepository[SavedUser](db)}
}

func (r *Repository) Get(ctx context.Context, id int64) (*user.User, error) {
	saved, err := r.engine.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return &saved.User, nil
}

func (r *Repository) List(ctx context.Context) ([]user.User, error) {
	saved, err := r.engine.List(ctx)
	if err != nil {
		return nil, err
	}
	models := make([]user.User, len(saved))
	for i := range saved {
		models[i] = saved[i].User
	}
	return models, nil
}

func (r *Repository) Create(ctx context.Context, model *user.User) error {
	saved := &SavedUser{User: *model}
	if err := r.engine.Create(ctx, saved); err != nil {
		return err
	}
	model.ID, model.UUID, model.CreatedAt, model.UpdatedAt = saved.ID, saved.UUID, saved.CreatedAt, saved.UpdatedAt
	return nil
}

func (r *Repository) Update(ctx context.Context, model *user.User) error {
	saved := &SavedUser{User: *model, ID: model.ID, UUID: model.UUID, CreatedAt: model.CreatedAt}
	if err := r.engine.Update(ctx, saved); err != nil {
		return err
	}
	model.UpdatedAt = saved.UpdatedAt
	return nil
}

func (r *Repository) Delete(ctx context.Context, model *user.User) error {
	return r.engine.Delete(ctx, &SavedUser{ID: model.ID})
}
```

`domain.User`'s own `ID`/`UUID`/`CreatedAt`/`UpdatedAt` fields are shadowed by
`SavedUser`'s tagged copies for bun's purposes — the adapter explicitly
copies values across at each boundary rather than relying on embedding
promotion. `application/user.Repository` and
`infrastructure/database/user.Repository` satisfy each other structurally
(Go interfaces); no explicit `var _ Repository = (*Repository)(nil)` wiring
is generated, though a consuming `main.go` can add one.

## Mechanism

- **`lib/generate/registry.go`**: `Kind` gains an `OutputPath` template
  (relative to project root, `{{.Package}}` substituted) replacing the flat
  `OutputFile`, since each kind now writes to a different directory shape
  and the application-service kind's filename itself is
  package-dependent (`{{.Package}}service.go`):
  ```go
  {Name: "model",       TemplateFile: "model.go.template",       OutputPath: "lib/domain/{{.Package}}/model.go"}
  {Name: "service",     TemplateFile: "service.go.template",     OutputPath: "lib/domain/{{.Package}}/service.go"}
  {Name: "app-service", TemplateFile: "appservice.go.template",  OutputPath: "lib/application/{{.Package}}/{{.Package}}service.go"}
  {Name: "repository",  TemplateFile: "repository.go.template",  OutputPath: "lib/infrastructure/database/{{.Package}}/repository.go"}
  ```
  Flags: `-m/--model`, `-s/--service`, `-a/--app-service`, `-r/--repository`.
- **`lib/generate/generate.go`**: `Run` substitutes `{{.Package}}` into each
  requested kind's `OutputPath` (plain string replace — package names are
  already-normalized lowercase identifiers, no template engine needed) to
  get the output path, instead of joining a single shared `lib/<pkg>/` dir
  with a flat filename.
- **`lib/cli/cli.go`**: unchanged — it already derives one subcommand per
  registry `Kind`, so the new `app-service` kind gets `letsgo make
  app-service <name>` for free.
- **`lib/templates/generators/`**: `appservice.go.template` added;
  `model.go.template`, `service.go.template`, `repository.go.template`
  rewritten to the shapes above.
- **`lib/templates/project/lib/database/`** moves to
  `lib/templates/project/lib/infrastructure/database/`. `model.go` (the
  shared `Model` base) is deleted — no longer needed, since UUID/timestamp
  handling now lives per-domain in each generated `Saved<Struct>`.
  `repository.go` and `connection.go` move unchanged.
  `cmd/http/main.go.template`'s import updates to
  `{{.AppName}}/lib/infrastructure/database`.

## Testing

- `lib/generate`: `registry_test.go` updated for the new `OutputPath` field
  and the 4-kind list (`model`, `service`, `app-service`, `repository`).
  `generate_test.go` updated to assert the new per-kind output locations and
  package/import contents.
- `lib/cli`: `cli_test.go` updated for the new output paths; add a case
  covering `letsgo make app-service <name>`.
- `lib/templates`: `templates_test.go` updated to check for
  `appservice.go.template` alongside the existing three.
- All four generated files remain asserted as syntactically valid Go via
  `go/parser`, as today.

## Extensibility

Unchanged: a future kind is one template file plus one registry entry.

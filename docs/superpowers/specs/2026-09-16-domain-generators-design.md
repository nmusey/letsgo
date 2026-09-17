# Domain file generators (`letsgo make model|repository|service`) — Design

## Context

`letsgo` currently has one generator: `letsgo make project`, which scaffolds a
whole new Go project from an embedded template tree
(`lib/templates/project/`), substituting `{{ .AppName }}` via
`lib/files.TemplateParser`. The scaffolded project ships generic, reusable
bun-backed pieces in `lib/database`: `Model` (base struct with UUID +
timestamps) and `Repository[T]` (generic CRUD).

There is no generator yet for individual domain-level files inside an
already-scaffolded project. This design adds one, modeled loosely on
Laravel's `artisan make:*` commands but scoped to logic-level files only —
no controllers or other "outside world" concerns, which come later.

## Goals

- `letsgo make model <name>`, `letsgo make repository <name>`, and
  `letsgo make service <name>` each generate one prepopulated Go file for a
  domain, inside an existing letsgo-scaffolded project.
- Files are organized **by concern, not by layer**: everything about one
  domain lives together in `lib/<name>/`, not spread across
  `lib/models/`, `lib/repositories/`, etc.
- Flags let one invocation generate multiple related files at once, e.g.
  `letsgo make model User -r -s` creates model, repository, and service
  together.
- The mechanism is extensible: adding a future generator kind is a new
  template + one registry entry, not a restructure.

## Non-goals

- Controllers or any transport/outside-world layer — explicitly deferred.
- Generating multiple domains in one invocation.
- Modifying/regenerating existing files — a collision is an error, not a
  merge or overwrite.

## Commands & flags

```
letsgo make model <name>      [-r|--repository] [-s|--service]
letsgo make repository <name> [-m|--model]      [-s|--service]
letsgo make service <name>    [-m|--model]      [-r|--repository]
```

Each subcommand always generates its own kind, plus whichever companion
flags are passed. `letsgo make repository Post -s` writes
`lib/post/repository.go` and `lib/post/service.go`, not `model.go`.

This is driven by one shared registry (kind → template file, long/short
flag name, struct/embed requirements), not three hand-duplicated command
bodies — each subcommand is "my kind, plus whatever's flagged."

One domain name per invocation for v1.

## Name normalization

The `<name>` argument is free-form — `BlogPost`, `blog_post`, `blog-post`,
`blog post`, `blogPost` are all accepted. It's split into words on
separators (space, `_`, `-`) and case boundaries (lower→upper transitions),
then rejoined two ways:

- **StructName** — PascalCase: `BlogPost`
- **PackageName** — all-lowercase, no separators: `blogpost` (matches Go's
  package-naming convention, which avoids underscores and mixed case)

## Generated code

Templates embed the existing generic bun helpers rather than duplicating
CRUD logic per domain. No placeholder/TODO comments are included — the
generated files are minimal and ready to extend.

```go
// lib/user/model.go
package user

import "{{.ModuleName}}/lib/database"

type User struct {
	database.Model
}
```

```go
// lib/user/repository.go
package user

import (
	"github.com/uptrace/bun"

	"{{.ModuleName}}/lib/database"
)

type Repository = database.Repository[User]

func NewRepository(db *bun.DB) *Repository {
	return database.NewRepository[User](db)
}
```

```go
// lib/user/service.go
package user

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}
```

`{{.ModuleName}}` is resolved by parsing the `module` line of `go.mod` in
the current working directory. This doubles as the project-root check: if
no `go.mod` is found there, the command fails with a clear error telling
the user to run it from their project root.

## Mechanism

- **Templates**: new embedded dir
  `lib/templates/generators/{model,repository,service}.go.template`,
  registered in `lib/templates/templates.go` alongside the existing
  `project`/`test` embeds (`//go:embed all:generators`).
- **`lib/generate` (new package)**: name normalization (word-splitting +
  PascalCase/lowercase joins), `go.mod` module-name resolution, the kind
  registry (kind name → template path → output filename → flag names), and
  an orchestrating `Generate(dir, kinds []string, structName, packageName,
  moduleName string) error` that each CLI subcommand calls with its
  resolved kind set.
- **`lib/files`**: new `RenderTemplate(fsys fs.FS, templatePath, outPath
  string, substitutions map[string]string) error` — a single-file sibling
  to the existing tree-walking `TemplateParser`. Opens `outPath` with
  `O_CREATE|O_EXCL|O_WRONLY` so an existing file fails immediately
  (collision = error, not overwrite) rather than needing a separate
  stat-then-write race. The target directory is created with `MkdirAll`
  first (idempotent, so two kinds writing into the same new
  `lib/<name>/` folder in one invocation don't conflict).
- **`lib/cli/cli.go`**: three new subcommands under `make`, each defining
  its own long/short companion flags and calling into `lib/generate`.

## Error handling

- No `go.mod` in cwd → error, nothing written.
- A target file already exists → that one file is reported as an error and
  skipped; other files requested in the same invocation still get written;
  the command exits non-zero if anything was skipped.

## Testing

- `lib/generate`: table-driven tests for name normalization (the
  space/underscore/hyphen/camelCase input variants above) and for the
  kind registry's flag-to-kind-set resolution.
- `lib/files`: test `RenderTemplate` writes expected content, substitutes
  correctly, and returns an error without writing when the target already
  exists.
- `lib/cli`: a couple of end-to-end tests driving `letsgo make model
  <name> -r -s` (and the repository/service equivalents) against a temp
  dir with a fake `go.mod`, asserting the right files land in
  `lib/<name>/` with the right package/struct names substituted.

## Extensibility

The registry is a small data table (kind → template/output/flags), so a
future kind (e.g. a real "other logic-level filetype") is one new template
file plus one registry entry — no changes to the orchestration or CLI
wiring logic.

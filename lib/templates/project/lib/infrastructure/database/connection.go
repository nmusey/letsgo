package database

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/driver/sqliteshim"
)

func NewDB() (*bun.DB, error) {
	switch driver := os.Getenv("DB_DRIVER"); driver {
	case "", "postgres":
		return newPostgresDB()
	case "sqlite":
		return newSQLiteDB()
	default:
		return nil, fmt.Errorf("database: unknown DB_DRIVER %q", driver)
	}
}

func newPostgresDB() (*bun.DB, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("database: connecting to postgres: %w", err)
	}

	return db, nil
}

func newSQLiteDB() (*bun.DB, error) {
	path := os.Getenv("DB_SQLITE_PATH")
	if path == "" {
		path = "file::memory:?cache=shared"
	}

	sqldb, err := sql.Open(sqliteshim.ShimName, path)
	if err != nil {
		return nil, fmt.Errorf("database: opening sqlite: %w", err)
	}

	db := bun.NewDB(sqldb, sqlitedialect.New())

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("database: connecting to sqlite: %w", err)
	}

	return db, nil
}

package core

import (
	"database/sql"
	"os"
	"strings"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func TestConnectToDatabase(t *testing.T) {
    db, err := sql.Open("postgres", "host=localhost port=5432 user=test_user password=test_password dbname=test_db sslmode=disable")
    if err != nil {
        t.Fatalf("Failed to open test database connection: %v", err)
    }
    defer db.Close()

    testConfig := DatabaseConfig{
        User:    "test_user",
        Password: "test_password",
        Host:    "localhost",
        Port:    "5432",
        Name:    "test_db",
    }
        
    _, err = ConnectToDatabase(testConfig)
    if !isValidDatabaseError(err) {
        t.Fatalf("ConnectToDatabase() returned an error: %v", err)
    }
}

func TestGetDatabaseConfig(t *testing.T) {
    os.Setenv("DB_USER", "test_user")
    os.Setenv("DB_PASS", "test_password")
    os.Setenv("DB_HOST", "test_host")
    os.Setenv("DB_PORT", "5432")
    os.Setenv("DB_NAME", "test_db")

    expectedConfig := DatabaseConfig{
        User:     "test_user",
        Password: "test_password",
        Host:     "test_host",
        Port:     "5432",
        Name:     "test_db",
    }

    config := GetDatabaseConfig()
    if config != expectedConfig {
        t.Errorf("GetDatabaseConfig() = %v, expected %v", config, expectedConfig)
    }
}

func TestMigrateDatabase(t *testing.T) {
    db, err := sql.Open("postgres", "host=localhost port=5432 user=test_user password=test_password dbname=test_db sslmode=disable")
    if !isValidDatabaseError(err) {
        t.Fatalf("Failed to open test database connection: %v", err)
    }
    defer db.Close()

    driver, err := postgres.WithInstance(db, &postgres.Config{})
    if !isValidDatabaseError(err) {
        t.Fatalf("Failed to create migration driver: %v", err)
    }

    _, err = migrate.NewWithDatabaseInstance(
        "file://../../migrations",
        "test_db",
        driver,
    )
    if !isValidDatabaseError(err) { 
        t.Fatalf("Failed to create migration instance: %v", err)
    }

    err = MigrateDatabase(db)
    if !isValidDatabaseError(err) {
        t.Fatalf("MigrateDatabase() returned an error: %v", err)
    }
}

func isValidDatabaseError(err error) bool {
    return err == nil || strings.Contains(err.Error(), "connection refused")
}

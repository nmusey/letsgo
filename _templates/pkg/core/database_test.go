package core

import (
	"os"
	"testing"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func TestDefaultDatabaseConfig(t *testing.T) {
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

    config := GetDefaultDatabaseConfig()
    if config != expectedConfig {
        t.Errorf("GetDatabaseConfig() = %v, expected %v", config, expectedConfig)
    }
}

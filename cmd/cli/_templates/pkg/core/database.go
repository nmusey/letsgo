package core

import (
    "database/sql"
    "fmt"
    "os"

    "github.com/golang-migrate/migrate/v4"
    "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
    _ "github.com/lib/pq"
)

type DatabaseConfig struct {
    User     string
    Password string
    Host     string
    Port     string
    Name     string
}

func ConnectToDatabase(config DatabaseConfig) (*sql.DB, error) {
    connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", config.Host, config.Port, config.User, config.Password, config.Name)
    
    db, err := sql.Open("postgres", connectionString)
    if err != nil {
        return nil, err
    }

    if err := db.Ping(); err != nil {
        return nil, err
    }

    return db, nil
}

func GetDatabaseConfig() DatabaseConfig {
    return DatabaseConfig{
        User:     os.Getenv("DB_USER"),
        Password: os.Getenv("DB_PASS"),
        Host:     os.Getenv("DB_HOST"),
        Port:     os.Getenv("DB_PORT"),
        Name:     os.Getenv("DB_NAME"),
    }
}

func MigrateDatabase(db *sql.DB) error {
    driver, err := postgres.WithInstance(db, &postgres.Config{})
    if err != nil {
        return err
    }

    config := GetDatabaseConfig()
    m, err := migrate.NewWithDatabaseInstance(
        "file://migrations",
        config.Name,
        driver,
    )
    if err != nil {
        return err
    }

    if err := m.Up(); err != nil {
        return err
    }

    return nil
}

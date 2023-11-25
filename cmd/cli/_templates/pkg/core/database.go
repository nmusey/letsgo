package core

import (
    "os"
    "fmt"
    "database/sql"

    _ "github.com/lib/pq"
    "github.com/golang-migrate/migrate/v4"
    "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
)

func ConnectToDatabase() (*sql.DB, error) {
    db, err := sql.Open("postgres", createConnectionString())
    if err != nil {
        return nil, err
    }

    if err := db.Ping(); err != nil {
        return nil, err 
    }

    return db, nil
}

func createConnectionString() string {
    user := os.Getenv("DB_USER")
    password := os.Getenv("DB_PASS")
    host := os.Getenv("DB_HOST")
    port := os.Getenv("DB_PORT")
    dbname := os.Getenv("DB_NAME")

    return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
}

func MigrateDatabase(db *sql.DB) error {
    driver, err := postgres.WithInstance(db, &postgres.Config{})
    if err != nil {
        return err
    }

    m, err := migrate.NewWithDatabaseInstance(
        "file://migrations",
        os.Getenv("DB_NAME"),
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

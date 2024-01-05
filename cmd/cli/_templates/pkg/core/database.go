package core

import (
    "database/sql"
    "fmt"
    "os"

    "github.com/jmoiron/sqlx"
    "github.com/golang-migrate/migrate/v4"
    "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
    _ "github.com/lib/pq"
)

type Database struct {
    DB      *sqlx.DB
    sqlDB   *sql.DB
    Config  DatabaseConfig
}

type DatabaseConfig struct {
    User     string
    Password string
    Host     string
    Port     string
    Name     string
}

func GetDefaultDatabaseConfig() DatabaseConfig {
    return DatabaseConfig{
        User:     os.Getenv("DB_USER"),
        Password: os.Getenv("DB_PASS"),
        Host:     os.Getenv("DB_HOST"),
        Port:     os.Getenv("DB_PORT"),
        Name:     os.Getenv("DB_NAME"),
    }
}

func (db *Database) Connect() error {
    connectionString := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", 
        db.Config.Host, 
        db.Config.Port, 
        db.Config.User, 
        db.Config.Password, 
        db.Config.Name,
    )

    sqlConnection, err := sql.Open("postgres", connectionString)
    if err != nil {
        return err
    }

    connection, err := sqlx.Connect("postgres", connectionString)

    db.DB = connection
    db.sqlDB = sqlConnection
    return nil
}

func (db Database) Migrate() error {
    driver, err := postgres.WithInstance(db.sqlDB, &postgres.Config{})
    if err != nil {
        return err
    }

    config := GetDefaultDatabaseConfig()
    migrator, err := migrate.NewWithDatabaseInstance(
        "file://migrations",
        config.Name,
        driver,
    )

    if err != nil {
        return err
    }

    if err := migrator.Up(); err != nil {
        return err
    }

    return nil
}

func (db Database) Select(dest interface{}, query string, args ...interface{}) error {
    return db.DB.Select(dest, query, args...)
}

func (db Database) SelectOne(dest interface{}, query string, args ...interface{}) error {
    return db.DB.Get(dest, query, args...)
}

func (db Database) NamedExec(query string, arg interface{}) (sql.Result, error) {
    return db.DB.NamedExec(query, arg)
}

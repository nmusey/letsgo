package core

import (
    "database/sql"
    "fmt"
    "os"
    "time"

    "github.com/golang-migrate/migrate/v4"
    "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
    _ "github.com/lib/pq"
)

type Database struct {
    SqlDB   *sql.DB
    Config  DatabaseConfig
}

type DatabaseConfig struct {
    User     string
    Password string
    Host     string
    Port     string
    Name     string
}

func NewDatabaseConnection(config DatabaseConfig) Database {
    db := Database{
        Config: config,
    }

    BlockingBackoff(db.Connect, 5, 3 * time.Second)
    return db
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
        fmt.Println(err)
        return err
    }

    db.SqlDB = sqlConnection
    return nil
}

func (db Database) Migrate() error {
    driver, err := postgres.WithInstance(db.SqlDB, &postgres.Config{})
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

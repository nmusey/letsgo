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

type Database struct {
    DB      *sql.DB
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

    connection, err := sql.Open("postgres", connectionString)
    if err != nil {
        return err
    }

    if err := connection.Ping(); err != nil {
        return err
    }

    db.DB = connection
    return nil
}

func (db Database) Migrate() error {
    driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
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

func (db Database) ReadOne(model IModel, condition string, arguments ...any) error {
    query := fmt.Sprintf("SELECT %s FROM %s WHERE %s", model.Columns(), model.Table(), condition)
    row := db.DB.QueryRow(query, arguments...)
    if row.Err() != nil {
        return row.Err()
    }

    return row.Scan(model.Populate())
}

func (db Database) Read(model IModel, condition string, arguments ...any) error {
    query := fmt.Sprintf("SELECT %s FROM %s %s", model.Columns(), model.Table(), condition)
    rows, err := db.DB.Query(query, arguments...)
	if err != nil {
		return err
	}

    defer rows.Close()
	for rows.Next() {
	}

	if err = rows.Err(); err != nil {
		return err
	}

    return nil
}

func (db Database) Insert(model IModel) error {
    query := fmt.Sprintf(
        "INSERT INTO %s (%s) VALUES (%s)", 
        model.Table(), 
        model.Columns(), 
        createColumnPlaceholders(model),
    )

    _, err := db.DB.Exec(query, model.ColumnValues()...)
    return err
}

func (db Database) Update(model IModel, query string, arguments ...any) error {
    _, err := db.DB.Exec(query, arguments...)
    return err
}

func (db Database) Delete(model IModel, condition string, arguments ...any) error {
    query := fmt.Sprintf("DELETE FROM %s WHERE %s", model.Table(), condition)
    _, err := db.DB.Exec(query, arguments...)
    return err
}

func createColumnPlaceholders(model IModel) string {
    columns := model.Columns()
    columnPlaceholders := ""
    for i := 0; i < len(columns); i++ {
        columnPlaceholders += fmt.Sprintf("$%d", i + 1)
        if i < len(columns) - 1 {
            columnPlaceholders += ", "
        }
    }

    return columnPlaceholders
}

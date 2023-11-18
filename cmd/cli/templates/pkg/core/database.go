package core

import (
    "os"
    "fmt"
    "database/sql"

    _ "github.com/lib/pq"
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

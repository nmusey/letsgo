package main

import (
	"fmt"
	"os"
	"time"

	"$appRepo/pkg/core"
)

func main() {
    db := core.Database{
        Config: core.DatabaseConfig{
            Host: os.Getenv("DB_HOST"),
            Port: os.Getenv("DB_PORT"),
            User: os.Getenv("DB_USER"),
            Password: os.Getenv("DB_PASS"),
            Name: os.Getenv("DB_NAME"),
        },
    }

    fmt.Println("Connecting to database...")
    core.BlockingBackoff(db.Connect, 5, 3 * time.Second)
    fmt.Println("Running migrations...")
    core.BlockingBackoff(db.Migrate, 5, 3 * time.Second)

    ctx := core.RouterContext{
        DB: db,
    }

    router := FiberRouter{
        ctx: &ctx,
    }

    router.Serve()
}

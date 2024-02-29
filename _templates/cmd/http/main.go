package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/handlebars/v2"

	"$appRepo/pkg/core"
)

func main() {
    app := fiber.New(fiber.Config{
        Views: handlebars.New("views", ".hbs.html"),
    })

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
        FiberApp: app,
    }

    router.RegisterRoutes()

    port := ":" + os.Getenv("APP_PORT")
    fmt.Printf("Listening on port %s\n", port)
    router.FiberApp.Listen(port)
}

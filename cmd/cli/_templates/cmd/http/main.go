package main

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/handlebars/v2"

	"$appRepo/pkg/core"
	"$appRepo/pkg/handlers"
	"$appRepo/pkg/middlewares/jwt"
)

func main() {
    app := fiber.New(fiber.Config{
        Views: handlebars.New("views", ".hbs.html"),
    })

    dbConfig := core.DatabaseConfig{
        Host: os.Getenv("DB_HOST"),
        Port: os.Getenv("DB_PORT"),
        User: os.Getenv("DB_USER"),
        Password: os.Getenv("DB_PASS"),
        Name: os.Getenv("DB_NAME"),
    }
    db, err := core.ConnectToDatabase(dbConfig)
    if err != nil {
        // Try again in 5 seconds, in case database is still booting up.
        time.Sleep(5 * time.Second)
        db, err = core.ConnectToDatabase(dbConfig)
        if err != nil {
            panic(err)
        }
    }

    if err := core.MigrateDatabase(db); err != nil {
        panic(err)
    }

    ctx := core.RouterContext{
        App: app,
        DB: db,
    }

    jwtSecret := os.Getenv("JWT_SECRET")
    app.Use(jwt.New(jwt.NewConfig(&ctx, jwtSecret)))

    handlers.NewUserHandler(&ctx).RegisterRoutes()
    handlers.NewAuthHandler(&ctx).RegisterRoutes()

    app.Listen(os.Getenv("APP_PORT"))
}

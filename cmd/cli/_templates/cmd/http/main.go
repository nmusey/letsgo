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

    db, err := core.ConnectToDatabase()
    if err != nil {
        // Try again in 5 seconds, in case database is still booting up.
        time.Sleep(5 * time.Second)
        db, err = core.ConnectToDatabase()
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

    app.Use(jwt.New(jwt.DefaultConfig(&ctx)))

    handlers.NewUserHandler(&ctx).RegisterRoutes()
    handlers.NewAuthHandler(&ctx).RegisterRoutes()

    app.Listen(os.Getenv("APP_PORT"))
}

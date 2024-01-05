package main

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/handlebars/v2"

	"$appRepo/pkg/core"
)

type FiberRouter struct {
    ctx *core.RouterContext
    FiberRouter *fiber.App
}

func (r FiberRouter) RegisterRoutes() {
    jwtSecret := os.Getenv("JWT_SECRET")
    r.FiberRouter.Use(jwt.New(jwt.NewConfig(r.ctx, jwtSecret)))

    routers := []core.Router{
        handlers.NewUserHandler(r.ctx, r.FiberRouter),
        handlers.NewAuthHandler(r.ctx, r.FiberRouter),
    }

    for _, router := range routers {
        router.RegisterRoutes()
    }
}

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

    core.BlockingBackoff(func() error {
        return db.Connect()
    }, 5, 3 * time.Second)

    if err := db.Migrate(); err != nil {
        panic(err)
    }

    ctx := core.RouterContext{
        DB: db,
    }

    router := FiberRouter{
        ctx: &ctx,
        FiberRouter: app,
    }

    router.RegisterRoutes()
    router.FiberRouter.Listen(os.Getenv("APP_PORT"))
}

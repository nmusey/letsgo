package main

import (
    "os"

	"github.com/gofiber/fiber/v2"

	"$appRepo/pkg/core"
	"$appRepo/pkg/handlers"
	"$appRepo/pkg/middlewares/jwt"
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

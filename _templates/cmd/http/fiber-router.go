package main

import (
	"os"

	"github.com/gofiber/fiber/v2"

	"$appRepo/pkg/core"
	"$appRepo/pkg/middlewares/jwt"
    "$appRepo/pkg/handlers"
)

type FiberRouter struct {
    ctx         *core.RouterContext
    FiberApp    *fiber.App
}

func (r FiberRouter) RegisterRoutes() {
    jwtSecret := os.Getenv("JWT_SECRET")
    jwtConfig := jwt.NewConfig(r.ctx, jwtSecret)
    r.FiberApp.Use(jwt.New(jwtConfig))

    r.RegisterAuthRoutes()
    r.RegisterUserRoutes()
}

func (r FiberRouter) Serve() {
    r.FiberApp = fiber.New()
    r.RegisterRoutes()
    r.FiberApp.Listen(":8080")
}

func (fr FiberRouter) RegisterAuthRoutes() {
    handler := handlers.NewAuthHandler(fr.ctx)

    fr.FiberApp.Get("/login", handler.GetLogin)
    fr.FiberApp.Get("/register", handler.GetRegister)
    fr.FiberApp.Get("/logout", handler.PostLogout)

    fr.FiberApp.Post("/login", handler.PostLogin)
    fr.FiberApp.Post("/register", handler.PostRegister)
    fr.FiberApp.Post("/logout", handler.PostLogout)
}

func (fr FiberRouter) RegisterUserRoutes() {
    handler := handlers.NewUserHandler(fr.ctx)

    fr.FiberApp.Get("/users", handler.GetUsers)
    fr.FiberApp.Get("/users/:id", handler.GetUserByID)
}

package main

import (
    "github.com/gofiber/fiber/v2"

    "$appRepo/pkg/core"
    "$appRepo/pkg/handlers"
)

func main() {
    app := fiber.New()
    db, err := core.ConnectToDatabase()
    if err != nil {
        panic(err)
    }

    ctx := core.RouterContext{
        App: app,
        DB: db,
    }


    handlers.NewUserHandler(&ctx).RegisterRoutes()
    
    app.Listen(":8080")
}

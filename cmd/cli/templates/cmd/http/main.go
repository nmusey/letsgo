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

    ctx := core.Context{
        DB: db,
    }


    users.NewUsersHandler(&ctx).RegisterRoutes()
    
    app.Listen(":8080")
}

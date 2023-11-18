package main

import (
    "github.com/gofiber/fiber/v2"

    "$appRepo/pkg/core"
    "$appRepo/pkg/users"
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


    users.NewHandler().RegisterRoutes(&ctx)
    
    app.Listen(":8080")
}

package main

import (
    "fmt"
    "net/http"

	"$appRepo/pkg/auth"
	"$appRepo/pkg/core"
	"$appRepo/pkg/users"
	"github.com/nmusey/letsgo/_templates/pkg/core"
)

func BuildRoutes(ctx *core.RouterContext) []core.Route {
    authHandler := auth.NewAuthHandler(ctx)
    userHandler := users.NewUserHandler(ctx)

    return []core.Route{
        core.BuildRoute("POST /migrate", handleMigrations)
        core.BuildRoute("GET /login", authHandler.GetLogin),
        core.BuildRoute("GET /register", authHandler.GetRegister),

        core.BuildRoute("POST /login", authHandler.PostLogin),
        core.BuildRoute("POST /register", authHandler.PostRegister),
        core.BuildRoute("POST /logout", authHandler.PostLogout),

        core.BuildRoute("GET /users", userHandler.GetUsers, &auth.JwtMiddleware{}, &core.JsonMiddleware{}),
    }
}

func handleMigrations(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Connecting to database...")
    core.BlockingBackoff(db.Connect, 5, 3 * time.Second)

    fmt.Println("Running migrations...")
    core.BlockingBackoff(db.Migrate, 5, 3 * time.Second)
}

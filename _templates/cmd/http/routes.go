package main

import (
	"$appRepo/pkg/auth"
	"$appRepo/pkg/core"
	"$appRepo/pkg/users"
)

func BuildRoutes(ctx *core.RouterContext) []core.Route {
    authHandler := auth.NewAuthHandler(ctx)
    userHandler := users.NewUserHandler(ctx)

    return []core.Route{
        core.BuildRoute("GET /login", authHandler.GetLogin),
        core.BuildRoute("GET /register", authHandler.GetRegister),

        core.BuildRoute("POST /login", authHandler.PostLogin),
        core.BuildRoute("POST /register", authHandler.PostRegister),
        core.BuildRoute("POST /logout", authHandler.PostLogout),

        core.BuildRoute("GET /users", userHandler.GetUsers, &auth.JwtMiddleware{}),
    }
}

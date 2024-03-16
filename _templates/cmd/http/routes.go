package main

import (
	"$appRepo/pkg/core"
	"$appRepo/pkg/handlers"
)

func BuildRoutes(ctx *core.RouterContext) map[string]core.HttpHandler {
    // TODO: Add jwt middleware here
    authHandler := handlers.NewAuthHandler(ctx)

    return map[string]core.HttpHandler{
        "GET /login": authHandler.GetLogin,
        "GET /register": authHandler.GetRegister,

        "POST /login": authHandler.PostLogin,
        "POST /register": authHandler.PostRegister,
        "POST /logout": authHandler.PostLogout,
    }
}

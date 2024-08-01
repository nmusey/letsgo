package main

import (
	"fmt"
	"net/http"
	"time"

	"$appRepo/internal/auth"
	"$appRepo/internal/core"
	"$appRepo/internal/users"
)

func BuildRoutes(router *core.Router) {
	authHandler := auth.NewAuthHandler(*router)
	userHandler := users.NewUserHandler(*router)

	router.Routes = []core.Route{
		core.BuildRoute("POST /migrate", buildMigrationHandler(*router)),

		core.BuildRoute("GET /login", authHandler.GetLogin),
		core.BuildRoute("GET /register", authHandler.GetRegister),

		core.BuildRoute("POST /login", authHandler.PostLogin),
		core.BuildRoute("POST /register", authHandler.PostRegister),
		core.BuildRoute("POST /logout", authHandler.PostLogout),

		core.BuildRoute("GET /users", userHandler.GetUsers, &auth.JwtMiddleware{}, &core.JsonMiddleware{}),
	}
}

func buildMigrationHandler(router core.Router) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Running migrations...")
		core.BlockingBackoff(router.DB.Migrate, 5, 3 * time.Second)
	}
}

package routes

import (
	"fmt"
	"net/http"
	"time"

	mw "$appRepo/internal/middleware"
	"$appRepo/internal/auth"
	"$appRepo/internal/core"
	"$appRepo/internal/users"
)

func BuildRoutes(router *core.Router) {
	authHandler := auth.NewAuthHandler(router)
	userHandler := users.NewUserHandler(router)

	router.Routes = []core.Route{
		core.BuildRoute("POST /migrate", buildMigrationHandler(*router)),

		core.BuildRoute("POST /login", authHandler.PostLogin),
		core.BuildRoute("POST /register", authHandler.PostRegister),
		core.BuildRoute("POST /logout", authHandler.PostLogout),

		core.BuildRoute("GET /users", userHandler.GetUsers, &mw.JwtMiddleware{}, &mw.JsonMiddleware{}),
        core.BuildRoute("GET /ping", handlePing),
        core.BuildRoute("GET /cache", buildCacheHandler(*router)),
	}
}

func handlePing(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("pong"))
}

func buildCacheHandler(router core.Router) http.HandlerFunc {
    return func (w http.ResponseWriter, r *http.Request) {
        router.Cache.Set("test", []byte("Hello world!"))
        resp, err := router.Cache.Get("test")
        if err != nil {
            w.WriteHeader(http.StatusInternalServerError)
        }

        w.Write(resp)
    }
}

func buildMigrationHandler(router core.Router) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Running migrations...")
		core.BlockingBackoff(router.DB.Migrate, 5, 3 * time.Second)
		fmt.Println("Migrations run")
	}
}

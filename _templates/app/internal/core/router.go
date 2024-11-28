package core

import (
	"fmt"
	"net/http"
	"os"
)

type Router struct {
    DB          Database
    Routes      []Route
    Mux         *http.ServeMux
    Cache       Cache
}

type Route struct {
    path        string
    handler     http.HandlerFunc
    middleware  []http.Handler
}

func (r *Router) Serve() {
    r.MapRoutes()

    port := ":" + os.Getenv("APP_PORT")
    fmt.Printf("Listening on port %s\n", port)
    http.ListenAndServe(port, r.Mux)
}

func BuildRoute(path string, handler http.HandlerFunc, middleware ...http.Handler) Route {
    return Route{
        path: path,
        handler: handler, 
        middleware: middleware,
    }
}

func (r *Router) MapRoutes() {
    for _, route := range r.Routes {
        r.Mux.HandleFunc(route.path, func(w http.ResponseWriter, r *http.Request) {
            for _, middleware := range route.middleware {
                middleware.ServeHTTP(w, r)
            }

            route.handler(w, r)
        })
    }
}

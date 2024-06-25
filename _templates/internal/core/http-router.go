package core

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/a-h/templ"

	"$appRepo/views/layouts"
)

type RouterContext struct {
    DB  Database
}

type Route struct {
    path        string
    handler     http.HandlerFunc
    middleware  []http.Handler
}

type HttpRouter struct {
    Mux *http.ServeMux
    ctx *RouterContext
}

func NewHttpRouter(ctx *RouterContext) *HttpRouter {
    return &HttpRouter{
        Mux: http.NewServeMux(),
        ctx: ctx,
    }
}

func (r *HttpRouter) Serve() {
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

func (r *HttpRouter) MapRoutes(routes []Route) *HttpRouter {
    for _, route := range routes {
        r.Mux.HandleFunc(route.path, func(w http.ResponseWriter, r *http.Request) {
            for _, middleware := range route.middleware {
                middleware.ServeHTTP(w, r)
            }

            route.handler(w, r)
        })
    }

    return r
}

func RenderTemplate(w http.ResponseWriter, components templ.Component) {
    layouts.MainLayout(components).Render(context.Background(), w)
}

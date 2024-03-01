package main

import (
	"net/http"
	"os"

	"$appRepo/pkg/core"
	"$appRepo/pkg/handlers"
	"$appRepo/pkg/middlewares/jwt"
)

type HttpRouter struct {
    ctx *core.RouterContext
    mux *http.ServeMux
}

func (r *HttpRouter) RegisterRoutes() {
    // TODO: Add routes here
    // TODO: Add jwt middleware here
}

func (r *HttpRouter) Serve() {
    r.mux = http.NewServeMux()
    r.RegisterRoutes()

    port := ":" + os.Getenv("APP_PORT")
    fmt.Printf("Listening on port %s\n", port)
    http.ListenAndServe(port, r.mux)
}

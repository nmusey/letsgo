package main

import (
    "fmt"
	"net/http"
	"os"

	"$appRepo/pkg/core"
	// "$appRepo/pkg/handlers"
	// "$appRepo/pkg/middlewares/jwt"
)

type HttpRouter struct {
    ctx *core.RouterContext
    mux *http.ServeMux
}

func (r *HttpRouter) RegisterRoutes() {
    // TODO: Add routes here
    // TODO: Add jwt middleware here
    r.mux.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("hello world"))
    })
}

func (r *HttpRouter) Serve() {
    r.mux = http.NewServeMux()
    r.RegisterRoutes()

    port := ":" + os.Getenv("APP_PORT")
    fmt.Printf("Listening on port %s\n", port)
    http.ListenAndServe(port, r.mux)
}

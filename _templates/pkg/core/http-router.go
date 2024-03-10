package core

import (
    "context"
    "encoding/json"
    "fmt"
	"net/http"
	"os"

	"github.com/a-h/templ"

	"$appRepo/views/layouts"
)

type HttpHandler func(http.ResponseWriter, *http.Request)error

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

func (r *HttpRouter) MapRoutes(routes map[string]HttpHandler) *HttpRouter {
    for route, handler := range routes {
        r.Mux.HandleFunc(route, func (w http.ResponseWriter, r *http.Request) {
            if err := handler(w, r); err != nil {
                w.WriteHeader(http.StatusInternalServerError)
                return
            }
        })
    }

    return r
}

func WriteJSON(w http.ResponseWriter, payload interface{}) error {
    w.Header().Set("Content-Type", "application/json")
    return json.NewEncoder(w).Encode(payload)
}

func ReadJSON(r *http.Request, payload interface{}) error {
    return json.NewDecoder(r.Body).Decode(payload)
}

func RenderTemplate(w http.ResponseWriter, components templ.Component) {
    layouts.MainLayout(components).Render(context.Background(), w)
}

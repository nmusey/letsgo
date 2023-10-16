package http

import (
    "net/http"
)

type Handler struct {
    Path string
    Method string
    Handler func (ctx Ctx) err
    Middleware []func (ctx Ctx) err
}

func (h *Handler) AddMiddleware(middleware func (ctx Ctx) err) {
    h.Middleware = append(h.Middleware, middleware)
}

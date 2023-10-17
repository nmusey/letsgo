package http

type Handler struct {
    Path string
    Method string
    Callback func (ctx Ctx) error
    Middleware []func (ctx Ctx) error
}

func (h *Handler) AddMiddleware(middleware func (ctx Ctx) error) {
    h.Middleware = append(h.Middleware, middleware)
}

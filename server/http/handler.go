package http

import (
    "net/http"
)

type Handler struct {
    Path string
    Method string
    Handler func (ctx Ctx) err
    Validators []func (ctx Ctx) err
}

func (h *Handler) AddValidator() {
    h.validators = append(h.validators, validator)
}

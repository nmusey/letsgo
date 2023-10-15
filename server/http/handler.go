package handler

import (
    "net/http"
)

type Handler struct {
    handler func (http.ResponseWriter, *http.Request): err
}

func (h *Handler) RegisterRoute(): http.Handler, err {
    return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
        err := h.handler(w, r)
        if err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            w.Write()
        }
    }), nil
}

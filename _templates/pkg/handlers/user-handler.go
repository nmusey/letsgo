package handlers

import (
    "github.com/gofiber/fiber/v2"
    "$appRepo/pkg/core"
    "$appRepo/pkg/services"
)

type UserHandler struct {
    ctx     *core.RouterContext
    Service services.UserService
}

func NewUserHandler(ctx *core.RouterContext) *UserHandler {
    return &UserHandler{
        ctx: ctx,
        Service: services.NewUserService(ctx),
    }
}

func (h UserHandler) GetUsers(w httpResponseWriter, r *http.Request) error {
    users, err := h.Service.GetUsers()
    if err != nil {
        return err
    }

    // TODO: Write JSON response with users
    return nil
}

func (h UserHandler) GetUserByID(w httpResponseWriter, r *http.Request) error {
    // TODO: Populate from request params
    var id int

    user, err := h.Service.GetUserByID(id)
    if err != nil {
        return err
    }

    // TODO: Write JSON response with user
    return nil
}

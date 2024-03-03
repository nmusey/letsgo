package handlers

import (
	"net/http"
	"strconv"

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

func (h UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) error {
    users, err := h.Service.GetUsers()
    if err != nil {
        return err
    }

    core.WriteJSON(w, users)
    return nil
}

func (h UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) error {
    var id int
    id, err := strconv.Atoi(r.PathValue("id"))

    user, err := h.Service.GetUserByID(id)
    if err != nil {
        return err
    }

    core.WriteJSON(w, user)
    return nil
}

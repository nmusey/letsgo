package handlers

import (
	"net/http"
	"strconv"

	"$appRepo/pkg/core"
	"$appRepo/pkg/services"
)

type UserHandler struct {
    ctx     *core.RouterContext
    UserService services.UserService
}

func NewUserHandler(ctx *core.RouterContext) *UserHandler {
    return &UserHandler{
        ctx: ctx,
        UserService: services.NewUserService(ctx),
    }
}

func (h UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) error {
    users, err := h.UserService.GetUsers()
    if err != nil {
        return err
    }

    core.WriteJSON(w, users)
    return nil
}

func (h UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) error {
    var id int
    id, err := strconv.Atoi(r.PathValue("id"))

    user, err := h.UserService.GetUserByID(id)
    if err != nil {
        return err
    }

    core.WriteJSON(w, user)
    return nil
}

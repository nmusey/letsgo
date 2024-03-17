package users 

import (
	"net/http"
	"strconv"

	"$appRepo/pkg/core"
)

type UserHandler struct {
    ctx     *core.RouterContext
    UserService UserService
}

func NewUserHandler(ctx *core.RouterContext) *UserHandler {
    return &UserHandler{
        ctx: ctx,
        UserService: NewUserService(ctx),
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

func (h UserHandler) GetUserById(w http.ResponseWriter, r *http.Request) error {
    var id int
    id, err := strconv.Atoi(r.PathValue("id"))

    user, err := h.UserService.GetUserById(id)
    if err != nil {
        return err
    }

    core.WriteJSON(w, user)
    return nil
}

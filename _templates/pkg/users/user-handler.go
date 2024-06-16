package users

import (
	"encoding/json"
	"net/http"
	"strconv"

	"$appRepo/pkg/core"
)

type userService interface {
    GetUsers() ([]User, error)
    GetUserById(int) (*User, error)
}

type UserHandler struct {
    ctx     *core.RouterContext
    UserService userService
}

func NewUserHandler(ctx *core.RouterContext) *UserHandler {
    return &UserHandler{
        ctx: ctx,
        UserService: NewUserService(ctx),
    }
}

func (h UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
    users, err := h.UserService.GetUsers()
    if err != nil {
        w.WriteHeader(http.StatusNotFound)
        return
    }

    response, err := json.Marshal(users)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    w.Write(response)
}

func (h UserHandler) GetUserById(w http.ResponseWriter, r *http.Request) {
    var id int
    id, err := strconv.Atoi(r.PathValue("id"))

    user, err := h.UserService.GetUserById(id)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    if user == nil {
        w.WriteHeader(http.StatusNotFound)
        return
    }

    response, err := json.Marshal(user)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    w.Write(response)
}

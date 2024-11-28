package users

import (
	"encoding/json"
	"net/http"
	"strconv"

	"$appRepo/internal/core"
)

type UserHandler struct {
    UserService *UserService
}

func NewUserHandler(router *core.Router) *UserHandler {
    return &UserHandler{
        UserService: NewUserService(router),
    }
}

func (h UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
    users, err := h.UserService.Store.GetUsers()
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

    user, err := h.UserService.Store.GetUserById(id)
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

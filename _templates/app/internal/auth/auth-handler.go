package auth

import (
	"encoding/json"
	"net/http"
	"net/mail"

	"$appRepo/internal/core"
	"$appRepo/internal/users"
)

type AuthHandler struct {
	service *AuthService
}

func NewAuthHandler(router *core.Router) *AuthHandler {
	return &AuthHandler{
		service: NewAuthService(router),
	}
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h AuthHandler) PostLogin(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.error(w, err, http.StatusBadRequest)
		return
	}

	user, err := h.service.UserService.Store.GetUserByEmail(req.Email)
	if err != nil {
        h.error(w, UserNotFoundError, http.StatusNotFound)
        return
	}

	correct, err := h.service.CheckPassword(req.Password, user.Id)
	if err != nil || !correct {
        h.error(w, UserNotFoundError, http.StatusNotFound)
        return 
	}

	cookie, err := h.service.GetAuthCookie(user)
	if err != nil {
        h.error(w, UserNotFoundError, http.StatusNotFound)
        return 
	}

	http.SetCookie(w, cookie)
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h AuthHandler) PostRegister(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.error(w, err, http.StatusBadRequest)
		return
	}

    err = h.validatePostRegister(req)
    if err != nil {
        h.error(w, err, http.StatusBadRequest)
        return
    }

	user := users.User{Email: req.Email}
    err = h.service.UserService.Store.SaveUser(&user)
    if err != nil {
        h.error(w, err, http.StatusBadRequest)
        return
    } 

    password := Password{Password: req.Password, UserId: user.Id}
    err = h.service.Store.savePassword(&password)
    if err != nil {
        h.error(w, err, http.StatusInternalServerError)
        return
    }

    cookie, err := h.service.GetAuthCookie(&user)
    if err != nil {
        h.error(w, err, http.StatusInternalServerError)
        return
    }

    http.SetCookie(w, cookie)
}

func (h AuthHandler) PostLogout(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:   "Authorization",
		Value:  "",
		MaxAge: -1,
	}

	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func (h AuthHandler) error(w http.ResponseWriter, err error, status int) {
	http.Error(w, err.Error(), status)
}

func (h AuthHandler) validatePostRegister(req RegisterRequest) error {
    if len(req.Password) < PasswordMinLength {
        return PasswordTooShortError
    }

    _, err := mail.ParseAddress(req.Email)
    if err != nil {
        return InvalidEmailError
    }

    return nil
}


package handlers

import (
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"$appRepo/pkg/core"
	"$appRepo/pkg/models"
	"$appRepo/pkg/services"
	"$appRepo/views/pages"
)

type AuthHandler struct {
    ctx             *core.RouterContext
    UserService     services.UserService
    PasswordService services.PasswordService
}

func NewAuthHandler(ctx *core.RouterContext) *AuthHandler {
    return &AuthHandler{
        ctx: ctx,
        UserService: services.NewUserService(ctx),
    }
}

func (h AuthHandler) GetLogin(w http.ResponseWriter, r *http.Request) error {
    core.RenderTemplate(w, pages.Login())
    return nil
}

func (h AuthHandler) GetRegister(w http.ResponseWriter, r *http.Request) error {
    core.RenderTemplate(w, pages.Register())
    return nil
}

func (h AuthHandler) PostLogin(w http.ResponseWriter, r *http.Request) error {
    r.ParseForm()
    email := r.FormValue("email")
    password := r.FormValue("password")

    user, err := h.UserService.GetUserByEmail(email)
    if err != nil {
        w.WriteHeader(http.StatusUnauthorized)
        return nil
    }

    match, err := h.PasswordService.CheckPassword(password, user.ID); 
    if !match || err != nil {
        w.WriteHeader(http.StatusUnauthorized)
        return nil
    }

    h.injectJwt(w, user)
    http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
    return nil
}

func (h AuthHandler) PostRegister(w http.ResponseWriter, r *http.Request) error {
    r.ParseForm()
    user := &models.User{
        Email: r.FormValue("email"),
    }

    if err := h.UserService.SaveUser(user); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return nil
    }

    user, err := h.UserService.GetUserByEmail(user.Email)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return nil
    }
    
    password := r.FormValue("password")
    if err := h.PasswordService.SavePassword(password, user.ID); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return nil
    }

    h.injectJwt(w, user)
    http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
    return nil
}

func (h AuthHandler) PostLogout(w http.ResponseWriter, r *http.Request) error {
    cookie := &http.Cookie{
        Name: "Authorization",
        Value: "",
        MaxAge: -1,
    }

    http.SetCookie(w, cookie)
    http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
    return nil
}

func (h AuthHandler) injectJwt(w http.ResponseWriter, user *models.User) error {
    expiry := time.Now().Add(time.Hour * 24)

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "uid": user.ID,
        "exp": expiry.Unix(),
    })

    tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
    if err != nil {
        return err
    }

    cookie := &http.Cookie{
        Name: "Authorization",
        Value: tokenString,
        Expires: expiry,
    }

    http.SetCookie(w, cookie)
    return nil
}

package handlers

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt"

	"$appRepo/pkg/core"
	"$appRepo/pkg/models"
	"$appRepo/pkg/services"
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
    return core.RenderTemplate(w, "pages/login", nil)
}

func (h AuthHandler) GetRegister(w http.ResponseWriter, r *http.Request) error {
    return core.RenderTemplate(w, "pages/register", nil)
}

func (h AuthHandler) PostLogin(w http.ResponseWriter, r *http.Request) error {
    r.ParseForm()
    email := r.FormValue("email")
    password := r.FormValue("password")

    user, err := h.UserService.GetUserByEmail(email)
    if err != nil {
        return err
    }

    match, err := h.PasswordService.CheckPassword(password, user.ID); 
    if !match || err != nil {
        return errors.New("Invalid email or password")
    }

    h.injectJwt(w, user)
    return w.Redirect("/")
}

func (h AuthHandler) PostRegister(w http.ResponseWriter, r *http.Request) error {
    r.ParseForm()
    user := models.User{
        Email: r.FormValue("email"),
        Username: r.FormValue("username"),
    }

    if err := h.UserService.SaveUser(user); err != nil {
        return err
    }

    user, err := h.UserService.GetUserByEmail(user.Email)
    if err != nil {
        return err
    }
    
    password := r.FormValue("password")
    if err := h.PasswordService.SavePassword(password, user.ID); err != nil {
        return err
    }

    h.injectJwt(w, user)
    return w.Redirect("/")
}

func (h AuthHandler) PostLogout(w http.ResponseWriter, r *http.Request) error {
    cookie := &http.Cookie{
        Name: "jwt",
        Value: "",
        MaxAge: -1,
    }

    http.SetCookie(w, cookie)
    return w.Redirect("/")
}

func (h AuthHandler) injectJwt(w http.ResponseWriter, user models.User) error {
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
        Name: "jwt",
        Value: tokenString,
        Expires: expiry,
    }

    http.SetCookie(w, cookie)
    return nil
}

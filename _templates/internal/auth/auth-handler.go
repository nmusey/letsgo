package auth

import (
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"$appRepo/internal/core"
	"$appRepo/internal/users"
	"$appRepo/views/pages"
)

type passwordService interface {
    SavePassword(string, int) error
    CheckPassword(string, int) (bool, error)
}

type userService interface {
    SaveUser(*users.User) error
    GetUserByEmail(string) (*users.User, error)
}

type AuthHandler struct {
    ctx             *core.RouterContext
    UserService     userService
    PasswordService passwordService
}

func NewAuthHandler(ctx *core.RouterContext) *AuthHandler {
    return &AuthHandler{
        ctx: ctx,
        UserService: users.NewUserService(ctx),
        PasswordService: NewPasswordService(ctx),
    }
}

func (h AuthHandler) GetLogin(w http.ResponseWriter, r *http.Request) {
    core.RenderTemplate(w, pages.Login())
}

func (h AuthHandler) GetRegister(w http.ResponseWriter, r *http.Request) {
    core.RenderTemplate(w, pages.Register())
}

func (h AuthHandler) PostLogin(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    email := r.FormValue("email")
    password := r.FormValue("password")

    user, err := h.UserService.GetUserByEmail(email)
    if err != nil {
        w.WriteHeader(http.StatusUnauthorized)
    }

    match, err := h.PasswordService.CheckPassword(password, user.Id); 
    if !match || err != nil {
        w.WriteHeader(http.StatusUnauthorized)
    }

    h.injectJwt(w, user)
    http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func (h AuthHandler) PostRegister(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    user := &users.User{
        Email: r.FormValue("email"),
    }

    if err := h.UserService.SaveUser(user); err != nil {
        w.WriteHeader(http.StatusBadRequest)
    }

    user, err := h.UserService.GetUserByEmail(user.Email)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
    }
    
    password := r.FormValue("password")
    if err := h.PasswordService.SavePassword(password, user.Id); err != nil {
        w.WriteHeader(http.StatusBadRequest)
    }

    h.injectJwt(w, user)
    http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func (h AuthHandler) PostLogout(w http.ResponseWriter, r *http.Request) {
    cookie := &http.Cookie{
        Name: "Authorization",
        Value: "",
        MaxAge: -1,
    }

    http.SetCookie(w, cookie)
    http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func (h AuthHandler) injectJwt(w http.ResponseWriter, user *users.User) {
    expiry := time.Now().Add(time.Hour * 24)
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "uid": user.Id,
        "exp": expiry.Unix(),
    })

    tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
    if err != nil {
        w.WriteHeader(http.StatusUnauthorized)
    }

    cookie := &http.Cookie{
        Name: "Authorization",
        Value: tokenString,
        Expires: expiry,
    }

    http.SetCookie(w, cookie)
}

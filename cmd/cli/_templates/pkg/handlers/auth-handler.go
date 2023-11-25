package handlers

import (
    "errors"
    "os"
    "time"

    "github.com/gofiber/fiber/v2"
    "github.com/golang-jwt/jwt"
    "$appRepo/pkg/core"
    "$appRepo/pkg/models"
    "$appRepo/pkg/services"
)

type AuthHandler struct {
    ctx *core.RouterContext
    UserService services.UserService
    PasswordService services.PasswordService
}

func NewAuthHandler(ctx *core.RouterContext) *AuthHandler {
    return &AuthHandler{
        ctx: ctx,
        UserService: services.NewUserService(ctx),
    }
}

func (h AuthHandler) RegisterRoutes() {
    h.ctx.App.Get("/login", h.GetLogin)
    h.ctx.App.Get("/register", h.GetRegister)
    h.ctx.App.Get("/logout", h.PostLogout)

    h.ctx.App.Post("/login", h.PostLogin)
    h.ctx.App.Post("/register", h.PostRegister)
    h.ctx.App.Post("/logout", h.PostLogout)
}

func (h AuthHandler) GetLogin(c *fiber.Ctx) error {
    return core.RenderTemplate("pages/login", c, fiber.Map{})
}

func (h AuthHandler) GetRegister(c *fiber.Ctx) error {
    return core.RenderTemplate("pages/register", c, fiber.Map{})
}

func (h AuthHandler) PostLogin(c *fiber.Ctx) error {
    email := c.FormValue("email")
    password := c.FormValue("password")

    user, err := h.UserService.GetUserByEmail(email)
    if err != nil {
        return err
    }

    match, err := h.PasswordService.CheckPassword(password, user.ID); 
    if !match || err != nil {
        return errors.New("Invalid email or password")
    }

    h.injectJwt(c, user)
    return c.Redirect("/", fiber.StatusOK)
}

func (h AuthHandler) PostRegister(c *fiber.Ctx) error {
    user := models.User{
        Email: c.FormValue("email"),
        Username: c.FormValue("username"),
    }

    if err := h.UserService.SaveUser(user); err != nil {
        return err
    }

    user, err := h.UserService.GetUserByEmail(user.Email)
    if err != nil {
        return err
    }
    
    password := c.FormValue("password")
    if err := h.PasswordService.SavePassword(password, user.ID); err != nil {
        return err
    }

    h.injectJwt(c, user)
    return c.Redirect("/")
}

func (h AuthHandler) PostLogout(c *fiber.Ctx) error {
    c.ClearCookie("jwt")
    return c.Redirect("/")
}

func (h AuthHandler) injectJwt(ctx *fiber.Ctx, user models.User) error {
    expiry := time.Now().Add(time.Hour * 24)

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "uid": user.ID,
        "exp": expiry.Unix(),
    })

    tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
    if err != nil {
        return err
    }

    ctx.Cookie(&fiber.Cookie{
        Name: "jwt",
        Value: tokenString,
        Expires: expiry,
    })

    return nil
}

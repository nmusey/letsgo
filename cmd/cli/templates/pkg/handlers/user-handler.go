package handlers

import (
    "github.com/gofiber/fiber/v2"
    "$appRepo/pkg/core"
    "$appRepo/pkg/models"
    "$appRepo/pkg/services"
)

type UserHandler struct {
    Service *services.UserService
}

func NewHandler(ctx *core.RouterContext) *UserHandler {
    return &UserHandler{
        Service: services.NewService(ctx),
    }
}

func (h UserHandler) RegisterRoutes(app *fiber.App) {
    app.Post("/users", h.SaveUser)
    app.Get("/users", h.GetUsers)
    app.Get("/users/:id", h.GetUserByID)
}

func (h UserHandler) SaveUser(c *fiber.Ctx) error {
    var user models.User
    if err := c.BodyParser(&user); err != nil {
        return err
    }

    err := h.Service.SaveUser(user)
    if err != nil {
        return err
    }

    return c.SendStatus(fiber.StatusCreated)
}

func (h UserHandler) GetUsers(c *fiber.Ctx) error {
    users, err := h.Service.GetUsers()
    if err != nil {
        return err
    }

    return c.JSON(users)
}

func (h UserHandler) GetUserByID(c *fiber.Ctx) error {
    id, err := c.ParamsInt("id")
    if err != nil {
        return err
    }

    user, err := h.Service.GetUserByID(id)
    if err != nil {
        return err
    }

    return c.JSON(user)
}

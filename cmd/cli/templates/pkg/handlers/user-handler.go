package handlers

import (
    "github.com/gofiber/fiber/v2"
    "$appRepo/pkg/core"
    "$appRepo/pkg/models"
    "$appRepo/pkg/services"
)

type UserHandler struct {
    ctx    *core.RouterContext
    Service *services.UserService
}

func NewUserHandler(ctx *core.RouterContext) *UserHandler {
    return &UserHandler{
        ctx: ctx,
        Service: services.NewService(ctx),
    }
}

func (h UserHandler) RegisterRoutes() {
    h.ctx.App.Post("/users", h.SaveUser)
    h.ctx.App.Get("/users", h.GetUsers)
    h.ctx.App.Get("/users/:id", h.GetUserByID)
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

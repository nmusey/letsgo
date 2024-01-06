package handlers

import (
    "github.com/gofiber/fiber/v2"
    "$appRepo/pkg/core"
    "$appRepo/pkg/services"
)

type UserHandler struct {
    ctx     *core.RouterContext
    Service services.UserService
}

func NewUserHandler(ctx *core.RouterContext) *UserHandler {
    return &UserHandler{
        ctx: ctx,
        Service: services.NewUserService(ctx),
    }
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

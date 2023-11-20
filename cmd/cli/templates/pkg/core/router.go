package core

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

type RouterContext struct {
    App *fiber.App
    DB *sql.DB
}

type Router interface {
    RegisterRoutes(ctx *RouterContext)
}

package core

import (
	"github.com/gofiber/fiber/v2"
)

func RenderTemplate(name string, ctx *fiber.Ctx, vars fiber.Map) error {
    return ctx.Render(name, vars, "layouts/main")
}

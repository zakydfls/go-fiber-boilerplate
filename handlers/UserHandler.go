package handlers

import (
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct{}

func (h *UserHandler) Test(ctx *fiber.Ctx) error {
	return ctx.SendString("Hello, Routes work!")
}

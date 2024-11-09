package handlers

import (
	"fiber_boilerplate/models"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct{}

var authModel = new(models.AuthModel)

func (a AuthHandler) Refresh(ctx *fiber.Ctx) {

}

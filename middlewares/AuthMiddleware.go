package middlewares

import (
	"fiber_boilerplate/models"

	"github.com/gofiber/fiber/v2"
)

var authModel = new(models.AuthModel)

func ValidateToken(ctx *fiber.Ctx) error {
	token, err := authModel.ExtractTokenMetadata(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Please login first",
		})
	}

	userId, err := authModel.GetAuth(token)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Please login first",
		})
	}
	ctx.Locals("userId", userId)
	return nil
}

package routes

import (
	"fiber_boilerplate/handlers"

	"github.com/gofiber/fiber/v2"
)

func InitRoutes(r *fiber.App) {
	api := r.Group("/api")

	v1 := api.Group("/v1")

	userHandler := new(handlers.UserHandler)
	v1.Get("/test", userHandler.Test)
}

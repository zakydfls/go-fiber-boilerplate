package routes

import (
	"fiber_boilerplate/handlers"

	"github.com/gofiber/fiber/v2"
)

func InitRoutes(r *fiber.App) {
	api := r.Group("/api")

	v1 := api.Group("/v1")

	auth := v1.Group("/auth")
	authHandler := new(handlers.AuthHandler)
	auth.Post("/register", authHandler.Register)
	auth.Post("/verify", authHandler.VerifyOtp)

	userHandler := new(handlers.UserHandler)
	v1.Get("/test", userHandler.Test)
}

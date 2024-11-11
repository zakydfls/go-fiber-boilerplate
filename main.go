package main

import (
	"fiber_boilerplate/db"
	"fiber_boilerplate/routes"
	validators "fiber_boilerplate/validator"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	validators.Initialize()

	app := fiber.New()
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))

	routes.InitRoutes(app)

	db.InitRedis(1)
	db.Init()

	app.Listen(":" + os.Getenv("PORT"))
}

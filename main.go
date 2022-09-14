package main

import (
	"os"

	v1 "github.com/Matterhorn-Studios/insiderviz-forms-api/v1"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	// init env
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	// create app
	app := fiber.New()

	// setup logger
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path} ${latency}\n",
	}))

	// add v1 group
	if err := v1.AddV1Group(app); err != nil {
		panic(err)
	}

	// add root ping
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})

	// setup port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// start server
	if err := app.Listen(":" + port); err != nil {
		panic(err)
	}
}

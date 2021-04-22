package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/pantame/server/config"
	"github.com/pantame/server/routes"
	"github.com/pantame/server/storage"
)

func main() {
	config.Read()
	storage.ConnectDB()
	storage.AutoMigrateDB()
	storage.ConnectCache()

	app := fiber.New(fiber.Config{
		//ProxyHeader: "X-Real-IP",
	})

	app.Use(cors.New())

	routes.Api1(app)

	app.Listen(":3000")
}

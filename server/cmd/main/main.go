package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/jrjaro18/tryingDC/internals/database"
	"github.com/jrjaro18/tryingDC/internals/routes"
	"github.com/jrjaro18/tryingDC/internals/rpc"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {

	go rpc.Init()

	app := fiber.New()
	app.Use(logger.New())
	//use cors
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders:  "Origin, Content-Type, Accept",
	}))

	api := app.Group("/api")
	routes.UserRoutes(api)
	routes.SellerRoutes(api)

	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("Pinging...")
	})

	// database.Connect()

	_, err := database.Init()
	if err != nil {
		panic(err)
	}

	app.Listen(":5000")
}

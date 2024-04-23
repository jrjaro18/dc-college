package main

import (
	"fmt"
	// "net"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/jrjaro18/tryingDC/internals/database"
	"github.com/jrjaro18/tryingDC/internals/redis"
	"github.com/jrjaro18/tryingDC/internals/routes"
	// "github.com/jrjaro18/tryingDC/internals/rpc"
)

func main() {
	//code is now saved separately in rpc folder outside of server folder
	// listener, err := net.Listen("tcp", ":0") // Listen on a random port
    // if err != nil {
    //     panic(err)
    // }

    // // Start the RPC server
	// go rpc.Init()

	go redis.Init()

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
	routes.ItemRoutes(api)

	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("Pinging...")
	})
	
	_, err := database.Init()
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to the database")

	app.Listen(":5000")
}

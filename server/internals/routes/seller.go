package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jrjaro18/tryingDC/internals/controllers"
)

func SellerRoutes(api fiber.Router) {
	route := api.Group("/seller")
	route.Post("/create", controllers.CreateSeller)
	route.Post("/login", controllers.LoginSeller)
	route.Post("/add-item", controllers.AddItem)
	route.Get("/get-logs", controllers.GetLogs)
}
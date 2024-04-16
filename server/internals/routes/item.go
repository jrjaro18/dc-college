package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jrjaro18/tryingDC/internals/controllers"
)

func ItemRoutes(app fiber.Router) {
	routes := app.Group("/item")
	routes.Get("/", controllers.GetAllItems)
}

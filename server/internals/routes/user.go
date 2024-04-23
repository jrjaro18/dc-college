package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jrjaro18/tryingDC/internals/controllers"
)

func UserRoutes(app fiber.Router) {
	routes := app.Group("/user")
	routes.Get("/", controllers.GetAllUsers)
	routes.Post("/create", controllers.CreateUser)
	routes.Post("/login", controllers.UserLogin)
	routes.Post("/add", controllers.AddToCart)
	routes.Post("/remove", controllers.RemoveFromCart)
	routes.Post("/cart", controllers.GetCart)
	routes.Post("/buy", controllers.Buy)
	routes.Post("/previouslyBought", controllers.PreviouslyBought)
}
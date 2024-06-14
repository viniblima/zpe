package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/viniblima/zpe/controllers"
	"github.com/viniblima/zpe/middlewares"
)

func SetupUserRoutes(api fiber.Router) {
	user_routes := api.Group("/users")

	user_routes.Post("/signup", controllers.SignUp)
	user_routes.Post("/signin", controllers.SignIn)

	user_routes.Get("/", middlewares.VerifyJWT, controllers.GetAllUsers)
	user_routes.Get("/:id", middlewares.VerifyJWT, controllers.GetUserDetail)
}

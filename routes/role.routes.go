package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/viniblima/zpe/controllers"
	"github.com/viniblima/zpe/middlewares"
)

func SetupRoleRoutes(api fiber.Router) {
	user_routes := api.Group("/roles")
	user_routes.Get("/", middlewares.VerifyJWT, controllers.GetRoles)

}

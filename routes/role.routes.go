package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/viniblima/zpe/controllers"
	"github.com/viniblima/zpe/middlewares"
)

/*
Configura as rotas de roles
*/
func SetupRoleRoutes(api fiber.Router) {
	user_routes := api.Group("/roles") // Configuracao da rota pai

	user_routes.Get("/", middlewares.VerifyJWT, controllers.GetRoles) // Listagem de roles criadas; necessário token de administrador ou modificador

}

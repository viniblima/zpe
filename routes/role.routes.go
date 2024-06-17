package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/viniblima/zpe/controllers"
	"github.com/viniblima/zpe/middlewares"
)

type RoleRouter interface {
	SetupRoleRoutes(api fiber.Router)
}
type roleController struct {
	roleControler controllers.RoleController
	middleware    middlewares.JWTMiddleware
}

/*
Configura as rotas de roles
*/
func (controller *roleController) SetupRoleRoutes(api fiber.Router) {
	user_routes := api.Group("/roles") // Configuracao da rota pai

	user_routes.Get("/", controller.middleware.VerifyJWT, controller.roleControler.GetAllRoles) // Listagem de roles criadas; necessário token de administrador ou modificador

}

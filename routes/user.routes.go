package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/viniblima/zpe/controllers"
	"github.com/viniblima/zpe/middlewares"
)

type UserRouter interface {
	SetupUserRoutes(api fiber.Router)
}
type userController struct {
	userController controllers.UserController
	middleware     middlewares.JWTMiddleware
}

/*
Configura as rotas de usuário
*/
func (controller *userController) SetupUserRoutes(api fiber.Router) {
	user_routes := api.Group("/users") // Configuracao da rota pai

	user_routes.Post("/signup", controller.userController.SignUp) // Criacao de usuário
	user_routes.Post("/signin", controller.userController.SignIn) // Login do usuário

	user_routes.Get("/", controller.middleware.VerifyJWT, controller.userController.GetAllUsers)      // Listagem de usuários; necessário token de visualizador
	user_routes.Get("/:id", controller.middleware.VerifyJWT, controller.userController.GetUserDetail) // Detalhe de usuário;  necessário token de visualizador

	user_routes.Patch("/:id", controller.middleware.VerifyJWTAdmin, controller.userController.UpdateUser) // Atualizacao de usuário; necessário token de administrador ou modificador

	user_routes.Delete("/:id", controller.middleware.VerifyJWTAdmin, controller.userController.RemoveUser) // Remove usuário
}

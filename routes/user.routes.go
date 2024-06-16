package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/viniblima/zpe/controllers"
	"github.com/viniblima/zpe/middlewares"
)

/*
Configura as rotas de usuário
*/
func SetupUserRoutes(api fiber.Router) {
	user_routes := api.Group("/users") // Configuracao da rota pai

	user_routes.Post("/signup", controllers.SignUp) // Criacao de usuário
	user_routes.Post("/signin", controllers.SignIn) // Login do usuário

	user_routes.Get("/", middlewares.VerifyJWT, controllers.GetAllUsers)      // Listagem de usuários; necessário token de visualizador
	user_routes.Get("/:id", middlewares.VerifyJWT, controllers.GetUserDetail) // Detalhe de usuário;  necessário token de visualizador

	user_routes.Patch("/:id", middlewares.VerifyJWTAdmin, controllers.UpdateUser) // Atualizacao de usuário; necessário token de administrador ou modificador
}

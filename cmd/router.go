package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/viniblima/zpe/routes"
)

type router struct {
	userRouter routes.UserRouter
	roleRouter routes.RoleRouter
}

/*
Esta funcao configura as rotas das APIs
*/
func (router *router) setupRoutes(app *fiber.App) {
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	api := app.Group("/api")

	// setupV1Routes(api)
	router.setupV1Routes(api)
}

/*
Esta funcao inicia as subrotas divida por versões
para facilitar manutenabilidade
e, na sequencia, divide por assuntos/tabelas
*/
func (router *router) setupV1Routes(api fiber.Router) {
	v1 := api.Group("/v1")

	router.userRouter.SetupUserRoutes(v1)
	router.roleRouter.SetupRoleRoutes(v1)
}

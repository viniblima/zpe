package controllers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	"github.com/viniblima/zpe/repository"
)

type RoleController interface {
	GetAllRoles(c *fiber.Ctx) error
}
type roleController struct {
	repo repository.RoleRepository
}

/*
Listagem de roles. A funcao se baseia nos seguintes passos:
1. Tentativa de trazer a listagem de todas as roles criadas; caso haja erro, retorna resposta com erro com o status 400
*/
func (controller *roleController) GetAllRoles(c *fiber.Ctx) error {

	roles, err := controller.repo.GetAllRoles()
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(err)
	}
	return c.Status(http.StatusOK).JSON(&roles)
}

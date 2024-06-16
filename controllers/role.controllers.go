package controllers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/viniblima/zpe/database"
	"github.com/viniblima/zpe/models"
)

/*
Listagem de roles. A funcao se baseia nos seguintes passos:
1. Tentativa de trazer a listagem de todas as roles criadas; caso haja erro, retorna resposta com erro com o status 400
*/
func GetRoles(c *fiber.Ctx) error {
	var roles []models.Role
	err := database.DB.Db.Model(models.Role{}).Find(&roles).Error
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(err)
	}
	return c.Status(http.StatusOK).JSON(&roles)
}

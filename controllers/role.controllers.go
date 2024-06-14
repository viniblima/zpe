package controllers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/viniblima/zpe/database"
	"github.com/viniblima/zpe/models"
)

func GetRoles(c *fiber.Ctx) error {
	var roles []models.Role
	database.DB.Db.Model(models.Role{}).Find(&roles)

	return c.Status(http.StatusOK).JSON(&roles)
}

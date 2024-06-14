package middlewares

import (
	"net/http"
	"os"
	"slices"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/viniblima/zpe/database"
	"github.com/viniblima/zpe/models"
)

func VerifyJWT(c *fiber.Ctx) error {

	auth := c.Get("Authorization")
	claims := jwt.MapClaims{}
	if auth != "" {
		split := strings.Split(auth, "JWT ")

		if len(split) < 2 {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid tag",
			})
		}

		_, err := jwt.ParseWithClaims(split[1], claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("PASSWORD_SECRET")), nil
		})

		if claims["sub"] != "auth" {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		if err != nil {
			return c.Status(401).JSON(fiber.Map{
				"error": err.Error(),
			})

		}

	} else {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid headers",
		})
	}

	c.Locals("userID", claims["id"])
	c.Next()
	return nil
}

func VerifyJWTAdmin(c *fiber.Ctx) error {

	auth := c.Get("Authorization")
	claims := jwt.MapClaims{}
	if auth != "" {
		split := strings.Split(auth, "JWT ")

		if len(split) < 2 {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid tag",
			})
		}

		_, err := jwt.ParseWithClaims(split[1], claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("PASSWORD_SECRET")), nil
		})

		if claims["sub"] != "auth" {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": err.Error(),
			})

		}

	} else {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid headers",
		})
	}

	c.Locals("userID", claims["ID"])

	var userAdmin models.User
	errUser := database.DB.Db.Model(models.User{}).Where("id = ?", claims["ID"]).Preload("Roles").First(&userAdmin).Error

	if errUser != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": errUser.Error(),
		})
	}

	indexRole := slices.IndexFunc(userAdmin.Roles, func(m *models.Role) bool {
		return m.Level == 1 || m.Level == 2
	})

	if indexRole < 0 {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "User doesnt have permission to do this",
		})
	}
	c.Next()
	return nil
}

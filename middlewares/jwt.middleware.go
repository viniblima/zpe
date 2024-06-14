package middlewares

import (
	"net/http"
	"os"
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
				"erro": "Invalid tag",
			})
		}

		_, err := jwt.ParseWithClaims(split[1], claims, func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodECDSA)
			if !ok {

			}
			return []byte(os.Getenv("PASSWORD_SECRET")), nil
		})

		if claims["sub"] != "auth" {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"erro": "Invalid token",
			})
		}

		if err != nil {
			return c.Status(401).JSON(fiber.Map{
				"erro": err.Error(),
			})

		}

	} else {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"erro": "Invalid headers",
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
				"erro": "Invalid tag",
			})
		}

		_, err := jwt.ParseWithClaims(split[1], claims, func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodECDSA)
			if !ok {

			}
			return []byte(os.Getenv("PASSWORD_SECRET")), nil
		})

		if claims["sub"] != "auth" {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"erro": "Invalid token",
			})
		}

		if err != nil {
			return c.Status(401).JSON(fiber.Map{
				"erro": err.Error(),
			})

		}

	} else {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"erro": "Invalid headers",
		})
	}

	c.Locals("userID", claims["id"])

	var userAdmin models.User
	database.DB.Db.Model(models.User{}).Where("ID = ?", claims["id"]).Where(models.Role{
		Level: 1,
	}).Association("Roles").Find(&userAdmin)
	c.Next()
	return nil
}

package middlewares

import (
	"net/http"
	"os"
	"slices"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/viniblima/zpe/models"
	"github.com/viniblima/zpe/repository"
)

type JWTMiddleware interface {
	VerifyJWT(c *fiber.Ctx) error
	VerifyJWTAdmin(c *fiber.Ctx) error
}
type jwtMiddleware struct {
	userRepo repository.UserRepository
	roleRepo repository.RoleRepository
}

/*
Funcao que verifica o JWT enviado pelo client. A funcao se baseia no seguintes passos:
 1. Verificacao se o cabecalho da requisicao está correto; caso haja erro, retorna resposta com erro com o status 401;
 2. Tentativa de extrair e validar o JWT enviado; caso haja erro, retorna resposta com erro com o status 401
 3. Insercao de ID do usuário nos Locals para acesso à esse dado posteriormente
    e para evitar que, mesmo que o JWT esteja correto,
    haja um contorno e tentativa de uma requisicao ese passando por outro usuário
*/
func (controller *jwtMiddleware) VerifyJWT(c *fiber.Ctx) error {

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

/*
Funcao que verifica o JWT enviado pelo client. A funcao se baseia no seguintes passos:
 1. Verificacao se o cabecalho da requisicao está correto; caso haja erro, retorna resposta com erro com o status 401;
 2. Tentativa de extrair e validar o JWT enviado; caso haja erro, retorna resposta com erro com o status 401
 3. Tentativa de verificacao das roles e níveis de permissao que o usuário possui para verrificar
    se pode ou nao efetuar aquela acao; caso o usuário não possua retorna resposta com erro com o status 401;
 4. Insercao de ID do usuário nos Locals para acesso à esse dado posteriormente
    e para evitar que, mesmo que o JWT esteja correto,
    haja um contorno e tentativa de uma requisicao ese passando por outro usuário
*/
func (controller *jwtMiddleware) VerifyJWTAdmin(c *fiber.Ctx) error {

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

	// var userAdmin models.User
	userAdmin, errUser := controller.userRepo.GetUserByID(c.Locals("userID").(string))
	// userAdmin, errUser := database.DB.Db.Model(models.User{}).Where("id = ?", claims["ID"]).Preload("Roles").First(&userAdmin).Error

	if errUser != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": errUser.Error(),
		})
	}

	indexRole := slices.IndexFunc(userAdmin.Roles, func(m *models.Role) bool {
		print(m.Level)
		return m.Level == 1 || m.Level == 2
	})

	if indexRole < 0 {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "User does not have permission to do this",
		})
	}

	roleLevel := slices.IndexFunc(userAdmin.Roles, func(m *models.Role) bool {
		return m.Level == 1
	})

	if roleLevel < 0 {
		c.Locals("roleLevel", 2)
	} else {
		c.Locals("roleLevel", 1)
	}

	c.Next()
	return nil
}

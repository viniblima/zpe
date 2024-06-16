package controllers

import (
	"errors"
	"net/http"
	"slices"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/viniblima/zpe/database"
	"github.com/viniblima/zpe/handlers"
	"github.com/viniblima/zpe/models"
	"gorm.io/gorm"
)

type LoginStruct struct {
	Email    string `json:"Email" validate:"required,email"`
	Password string `json:"Password" validate:"required"`
}

/*
Funcao que extrai os dados corretos à serem enviados ao client
*/
func extractUserObj(u models.User) map[string]interface{} {
	return map[string]interface{}{
		"ID":    u.ID,
		"Name":  u.Name,
		"Email": u.Email,
		"Roles": u.Roles,
	}
}

/*
Funcao que busca um usuário pelo email
*/
func getByEmail(email string) (models.User, error) {
	item := models.User{}
	var err error
	dbResult := database.DB.Db.Where("email = ?", email).Preload("Roles").First(&item)

	if errors.Is(dbResult.Error, gorm.ErrRecordNotFound) {
		err = errors.New("user not found")
	}

	return item, err
}

/*
Funcao de login do usuário, baseada nos seguintes passos:
1. Validacao do body enviado pela requisicao;
2. Tentativa de busca do usuário pelo email; caso haja erro, retorna resposta com erro com o status 401;
3. Verificacao de password no body com o passowrd criptografado no banco de dados; caso haja erro, retorna resposta com erro com o status 401;
4 Tentativa de geracao do token JWT para acesso às rotas protegidas; caso haja erro, retorna resposta com erro com o status 400
*/
func SignIn(c *fiber.Ctx) error {

	body := new(LoginStruct)

	c.BodyParser(&body)

	err := validator.New().Struct(body)
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(handlers.NewJError(err))
	}

	u, errorU := getByEmail(body.Email)

	checked := handlers.CheckHash(u.Password, body.Password)

	if errorU != nil || !checked {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": "Email or password wrong"})
	}

	json, errJwt := handlers.GenerateJWT(u.ID)

	if errJwt != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "Error on generate JWT"})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"Auth": json,
		"User": extractUserObj(u),
	})

}

type SignUpStruct struct {
	Email    string `json:"Email" validate:"required,email"`
	Password string `json:"Password" validate:"required"`
	Name     string `json:"Name" validate:"required"`
}

/*
Criacao de usuário com o body enviado. A funcao se baseia nos seguintes passos

1. Utilizacao de um struct para filtrar os dados necessários para criacao do usuário
2. Validacao dos campos
3. Tentativa de insercao dos dados no banco de dados
4. Geracao de JWT para autorizacao de acesso às rotas protegidas
5. Envio de JWT e dados do usuário
*/
func SignUp(c *fiber.Ctx) error {

	body := new(SignUpStruct)

	c.BodyParser(&body)

	err := validator.New().Struct(body)
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(handlers.NewJError(err))
	}
	user := models.User{
		Name:     body.Name,
		Password: body.Password,
		Email:    body.Email,
	}

	err = database.DB.Db.Create(&user).Error
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(handlers.NewJError(err))
	}
	json, errJwt := handlers.GenerateJWT(user.ID)

	if errJwt != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "Error on generate JWT"})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"Auth": json,
		"User": extractUserObj(user),
	})
}

/*
Listagem de usuários. A funcao se baseia nos seguintes passos:
1. Tentativa de trazer listagem de todos os usuários presentes, com a omissão do dado Password; caso haja erro, retorna resposta com erro com o status 400
*/
func GetAllUsers(c *fiber.Ctx) error {
	var users []models.User

	err := database.DB.Db.Omit("Password").Model(&models.User{}).Preload("Roles").Find(&users).Error
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(err)
	}
	return c.Status(http.StatusOK).JSON(users)

}

/*
Funcao que traz os detalhes do usuário, com os seguintes passos:
1. Tentativa de selecao do usuário no Banco de dados; caso haja erro, retorna resposta com erro com o status 400
*/
func GetUserDetail(c *fiber.Ctx) error {
	var users models.User

	id := c.Params("id")

	err := database.DB.Db.Omit("Password").Model(&models.User{}).Where("ID = ?", id).Preload("Roles").First(&users).Error
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "User not found",
		})
	}
	return c.Status(http.StatusOK).JSON(users)
}

type IDSRole struct {
	List []string `json:"List" validate:"required"`
}

/*
Funcao de atualizacao do dados do usuário, com os seguintes passos
1. Validacao da estrutura enviada no body da requisicao; ; caso haja erro, retorna resposta com erro com o status 422 com o detalhamento do erro
2. Busca pelo usuário com o ID enviado via parâmetro de URL
3. Busca pelas roles enviadas via body da requisicao;
4. Atualizacao das roles do usuário selecionadas pela listagem do body
*/
func UpdateUser(c *fiber.Ctx) error {
	body := new(IDSRole)

	c.BodyParser(&body)

	err := validator.New().Struct(body)
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(handlers.NewJError(err))
	}

	var users models.User

	id := c.Params("id")

	errUser := database.DB.Db.Model(&models.User{}).Where("ID = ?", id).First(&users).Error

	if errUser != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	var roles []models.Role

	roleLevel := c.Locals("roleLevel").(int)

	database.DB.Db.Where("ID IN ?", body.List).Find(&roles)

	indexRole := slices.IndexFunc(roles, func(m models.Role) bool {
		return m.Level < roleLevel
	})

	if indexRole >= 0 {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "User does not have the role necessary to give this role",
		})
	}

	database.DB.Db.Model(&users).Association("Roles").Clear()

	database.DB.Db.Model(&users).Omit("Roles.*").Association("Roles").Append(&roles)

	return c.Status(http.StatusOK).JSON(fiber.Map{

		"user": &users,
	})
}

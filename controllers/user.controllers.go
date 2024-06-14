package controllers

import (
	"errors"
	"net/http"

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

func extractUserObj(u models.User) map[string]interface{} {
	return map[string]interface{}{
		"ID":    u.ID,
		"Name":  u.Name,
		"Email": u.Email,
		"Roles": u.Roles,
	}
}

func getByEmail(email string) (models.User, error) {
	item := models.User{}
	var err error
	dbResult := database.DB.Db.Where("email = ?", email).Preload("Roles").First(&item)

	if errors.Is(dbResult.Error, gorm.ErrRecordNotFound) {
		err = errors.New("user not found")
	}

	return item, err
}

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

func GetAllUsers(c *fiber.Ctx) error {
	var users []models.User

	err := database.DB.Db.Omit("Password").Model(&models.User{}).Preload("Roles").Find(&users).Error
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(err)
	}
	return c.Status(http.StatusOK).JSON(users)

}

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

// type ListIDRoles struct {
// 	List []IDRole `validate:"required"`
// }

func UpdateUser(c *fiber.Ctx) error {
	body := new(IDSRole)

	c.BodyParser(&body)

	err := validator.New().Struct(body)
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(handlers.NewJError(err))
	}

	var users models.User

	id := c.Params("id")

	errUser := database.DB.Db.Omit("Password").Model(&models.User{}).Where("ID = ?", id).First(&users).Error

	if errUser != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	var roles []models.Role

	database.DB.Db.Where("ID IN ?", body.List).Find(&roles)

	return c.Status(http.StatusOK).JSON(&roles)
}

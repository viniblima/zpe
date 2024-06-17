package repository

import (
	"github.com/viniblima/zpe/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(todoItem *models.User) (*models.User, error)
	GetAllUsers() (*[]models.User, error)
	RemoveUser(user *models.User) error
	UpdateUser(user *models.User) (*models.User, error)
	UpdateUserRoles(user *models.User, roles *[]models.Role) *models.User
	GetUserByEmail(email string) (*models.User, error)
	GetUserDetail(id string) (*models.User, error)
	RemoveUserByID(id string) error
	GetUserByID(id string) (*models.User, error)
}

type userRepository struct {
	Db *gorm.DB
}

// Cria um novo usuário
func (repo *userRepository) CreateUser(user *models.User) (*models.User, error) {
	err := repo.Db.Create(user).Error
	return user, err
}

// Pega lista de usuários
func (repo *userRepository) GetAllUsers() (*[]models.User, error) {
	var items []models.User
	err := repo.Db.Find(&items).Error
	return &items, err
}

func (repo *userRepository) GetUserDetail(id string) (*models.User, error) {

	var user models.User
	// err := database.Db.Find(&items).Error
	err := repo.Db.Omit("Password").Model(&models.User{}).Where("ID = ?", id).Preload("Roles").First(&user).Error
	return &user, err
}

// Remove usuário
func (repo *userRepository) RemoveUser(user *models.User) error {
	err := repo.Db.Delete(user).Error
	return err
}

// Remove usuário por ID
func (repo *userRepository) RemoveUserByID(id string) error {
	var user models.User
	err := repo.Db.Model(&models.User{}).Where("ID = ?", id).Delete(user).Error
	return err
}

/*
Atualiza os dados de um usuário
*/
func (repo *userRepository) UpdateUser(user *models.User) (*models.User, error) {
	err := repo.Db.Save(user).Error
	return user, err
}

/*
Atualiza as roles de um usuário
*/
func (repo *userRepository) UpdateUserRoles(user *models.User, roles *[]models.Role) *models.User {
	repo.Db.Model(&user).Association("Roles").Clear()
	repo.Db.Model(&user).Omit("Roles.*").Association("Roles").Append(&roles)
	return user
}

/*
Funcao que busca um usuário pelo email
*/
func (repo *userRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := repo.Db.Where("email = ?", email).Preload("Roles").First(&user).Error

	return &user, err
}

/*
Funcao que busca um usuário pelo ID
*/
func (repo *userRepository) GetUserByID(id string) (*models.User, error) {
	var user models.User
	err := repo.Db.Where("ID = ?", id).Preload("Roles").First(&user).Error

	return &user, err
}

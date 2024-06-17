package repository

import (
	"github.com/viniblima/zpe/database"
	"github.com/viniblima/zpe/models"
	"gorm.io/gorm"
)

type RoleRepository interface {
	GetAllRoles() (*[]models.Role, error)
	GetRolesWithList(list []string) (*[]models.Role, error)
}

type roleRepository struct {
	Db *gorm.DB
}

// Pega lista de roles
func (repo *roleRepository) GetAllRoles() (*[]models.Role, error) {
	var items []models.Role
	err := database.Db.Find(&items).Error
	return &items, err
}

func (repo *roleRepository) GetRolesWithList(list []string) (*[]models.Role, error) {
	var items []models.Role
	// err := database.Db.Find(&items).Error
	err := database.Db.Where("ID IN ?", list).Find(&items).Error
	return &items, err
}

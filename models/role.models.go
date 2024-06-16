package models

import (
	uuid "github.com/satori/go.uuid"
	gorm "gorm.io/gorm"
)

/*
Modelo de Role
*/
type Role struct {
	gorm.Model
	ID    string  `gorm:"primaryKey"`
	Name  string  `json:"Name" validate:"required,min=3,max=32"`
	Level int     `gorm:"unique" json:"Level" validate:"required, min=1"`
	Users []*User `gorm:"many2many:user_roles;" json:"Users"`
}

func (m *Role) BeforeCreate(db *gorm.DB) (err error) {
	m.ID = uuid.NewV4().String()
	return
}

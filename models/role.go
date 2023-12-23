package models

import (
	"time"

	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	Id        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"string" gorm:"primaryKey"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func CreateRole(role *[]Role, db *gorm.DB) *gorm.DB {
	result := db.Create(role)
	return result
}

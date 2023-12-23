package models

import (
	"time"

	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	Id            uint      `json:"id" gorm:"primaryKey"`
	Message       string    `json:"token"`
	TypeOfMessage string    `json:"type_of_message"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
}

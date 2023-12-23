package models

import (
	"time"

	"gorm.io/gorm"
)

type Notification struct {
	gorm.Model
	Id         uint   `json:"id" gorm:"primaryKey"`
	Message    string `json:"message"`
	UserId     uint
	User       User      `gorm:"constraint:OnDelete:CASCADE"`
	ExtraValue string    `json:"extra_value"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
}

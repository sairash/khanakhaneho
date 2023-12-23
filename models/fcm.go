package models

import (
	"time"

	"gorm.io/gorm"
)

type Fcm struct {
	gorm.Model
	Id        uint   `json:"id" gorm:"primaryKey"`
	Token     string `json:"token"`
	UserId    uint
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func CreateFcm(fcm *Fcm, db *gorm.DB) *gorm.DB {
	return db.Create(fcm)
}

package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Id           uint   `json:"id" gorm:"primaryKey"`
	PhoneNumber  int    `json:"phone_number"  gorm:"index:unique"`
	Password     string `json:"password"`
	Token        string
	Notification []Notification
	Fcm          []Fcm
	MessageId    uint
	Message      Message
	RoleId       uint `gorm:"default:2"`
	Role         Role
	FriendId     int
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}

func CreateUsers(user *[]User, db *gorm.DB) {
	db.Create(&user)
}

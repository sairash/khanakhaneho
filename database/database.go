package database

import (
	"fmt"
	"khanakhaneho/helper"
	"khanakhaneho/models"
	"log"
	"strconv"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func ConnectDb() {
	db, err := gorm.Open(sqlite.Open("khanakhaneho.db"), &gorm.Config{})

	if err != nil {
		fmt.Println(err)
	}

	if err = db.AutoMigrate(&models.User{}, &models.Fcm{}, &models.Message{}, &models.Role{}, &models.Notification{}); err == nil {
		if (db.Migrator().HasTable(&models.Role{})) {
			if err = db.First(&models.Role{}).Error; err != nil {
				role := []models.Role{
					{Name: "Admin"},
					{Name: "Son"},
					{Name: "Mom"},
				}
				fmt.Println(models.CreateRole(&role, db))
			}
		}

		if (db.Migrator().HasTable(&models.User{})) {
			if err = db.First(&models.User{}).Error; err != nil {
				hashedPassword, err := bcrypt.GenerateFromPassword([]byte(helper.EnvVariable("ADMIN_PASSWORD")), bcrypt.DefaultCost)
				if err != nil {
					log.Fatal(err)
				}

				hashedUserPassword, err := bcrypt.GenerateFromPassword([]byte(helper.EnvVariable("USER_PASSWORD")), bcrypt.DefaultCost)
				if err != nil {
					log.Fatal(err)
				}

				phone_, err := strconv.Atoi(helper.EnvVariable("ADMIN_PHONE"))
				if err != nil {
					panic(err)
				}

				son, err := strconv.Atoi(helper.EnvVariable("SON_PHONE"))
				if err != nil {
					panic(err)
				}
				mom, err := strconv.Atoi(helper.EnvVariable("MOM_PHONE"))
				if err != nil {
					panic(err)
				}
				user := []models.User{
					{
						PhoneNumber: phone_,
						RoleId:      1,
						Token:       helper.RandomString(10),
						Password:    string(hashedPassword),
						FriendId:    0,
					},
					{
						PhoneNumber: son,
						RoleId:      2,
						Token:       helper.RandomString(10),
						Password:    string(hashedUserPassword),
						FriendId:    3,
					},
					{
						PhoneNumber: mom,
						RoleId:      3,
						Token:       helper.RandomString(10),
						Password:    string(hashedUserPassword),
						FriendId:    2,
					},
				}

				models.CreateUsers(&user, db)

			}
		}
	}

	helper.Database = helper.DbInstance{
		Db: db,
	}
}

package controller

import (
	"fmt"
	"khanakhaneho/helper"
	"khanakhaneho/models"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type (
	loginUser struct {
		PhoneNumber string `json:"phone_number"`
		Password    string `json:"password"`
		Fcm         string `json:"fcm"`
	}

	changeFriendStruct struct {
		ChangeToUser int  `json:"change_to_user"`
		Applying     bool `json:"applying"`
		Accepting    bool `json:"accepting"`
	}

	SendMessageStruct struct {
		Message       string `json:"message"`
		TypeOfMessage string `json:"type_of_message"`
	}

	CombinedUser struct {
		Id          uint
		PhoneNumber int
		RoleId      uint
		FriendPhone int
		FriendRole  uint
	}
)

func Me(c echo.Context) error {
	claims, err := helper.JWT(c)

	if err != nil {
		return helper.ErrorResponse(c, 2, err.Error(), nil)
	}

	var user CombinedUser

	// if err := helper.JWTAuthUser(claims.Token, &user); err != nil {
	// 	return helper.ErrorResponse(c, "Error occoured while fetching data.", nil)
	// }

	if err := helper.Database.Db.
		Select("users.id, users.phone_number, users.role_id, friends.role_id AS \"friend_role\", friends.phone_number AS \"friend_phone\"").
		Model(&models.User{}).
		Joins("JOIN users AS friends ON users.friend_id = friends.id").
		Where("users.token = ?", claims.Token).Scan(&user).Error; err != nil {
		return helper.ErrorResponse(c, 2, "Something Went Wrong!", nil)
	}
	if user.Id == 0 {
		only_user := models.User{}
		helper.Database.Db.Select("id", "role_id", "phone_number", "friend_id").Where("token = ?", claims.Token).First(&only_user)
		if only_user.Id == 0 {
			return helper.ErrorResponse(c, 1, "Error Authentication!", nil)
		}
		user.Id = only_user.Id
		user.PhoneNumber = only_user.PhoneNumber
		user.RoleId = only_user.RoleId
	}

	return helper.SuccessResponse(c, 0, "Token Success", echo.Map{
		"user": user,
	})

}

func Signup(password string) models.User {
	user := models.User{
		PhoneNumber: 9842519566,
		RoleId:      2,
		Token:       helper.RandomString(10),
		Password:    string(password),
		FriendId:    3,
	}
	helper.Database.Db.Create(&user)
	fmt.Println(user)
	return user
}

func Login(c echo.Context) error {
	request := loginUser{}
	if err := c.Bind(&request); err != nil {
		return helper.ErrorResponse(c, 1, "Error Validation", nil)
	}

	if request.PhoneNumber == "" || request.Password == "" || request.Fcm == "" {
		return helper.ErrorResponse(c, 1, "Please fill in something!", nil)
	}

	user := models.User{}

	helper.Database.Db.Select("id", "password", "RoleId", "token").Where("phone_number = ?", request.PhoneNumber).Find(&user)

	if user.Id == 0 {
		user = Signup(request.Password)
	} else {
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))

		if err != nil {
			return helper.ErrorResponse(c, 1, "Username or Password error!", nil)
		}
	}

	var fcm models.Fcm
	result := helper.Database.Db.FirstOrCreate(&fcm, models.Fcm{
		UserId: user.Id,
		Token:  request.Fcm,
	})

	if result.RowsAffected > 1 {
		// helper.Notification(fcm, )
		// send_fcm_notification(, "Friend Request", fmt.Sprintf("%d sent you a friend request.", user.PhoneNumber))

	}

	claims := &helper.UserJwtClaims{
		Token: user.Token,
		Role:  int(user.RoleId),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(helper.EnvVariable("SERECT")))
	if err != nil {
		return helper.ErrorResponse(c, 1, "Token Not Correct", nil)
	}

	return helper.SuccessResponse(c, 0, "Logged in successfully", echo.Map{
		"token":     t,
		"user_id":   user.Id,
		"user_type": user.RoleId,
	})

}

func SendMessage(c echo.Context) error {
	claims, err := helper.JWT(c)

	if err != nil {
		return helper.ErrorResponse(c, 2, err.Error(), nil)
	}

	request := SendMessageStruct{}
	if err := c.Bind(&request); err != nil {
		return helper.ErrorResponse(c, 1, "Error Validation", nil)
	}

	if request.Message == "" || request.TypeOfMessage == "" {
		return helper.ErrorResponse(c, 1, "Friend number is not correct", nil)
	}

	var user models.User
	if err := helper.Database.Db.First(&user, "token =?", claims.Token).Error; err != nil {
		return helper.ErrorResponse(c, 2, "Something Went Wrong!", nil)
	}

	new_message := models.Message{
		Message:       request.Message,
		TypeOfMessage: request.TypeOfMessage,
	}

	helper.Database.Db.Save(&new_message)

	var title string
	var body string

	if request.TypeOfMessage == "question" {
		title = "Question was asked"
		body = "Question: " + strings.Join(strings.Split(request.Message, "_"), " ")
	} else {
		title = "Answer was sent"
		body = "Answer: " + strings.Join(strings.Split(request.Message, "_"), " ")
	}

	send_fcm_notification(user.FriendId, title, body)

	helper.Database.Db.Model(&models.User{}).Where("id IN ?", []string{strconv.Itoa(int(user.Id)), strconv.Itoa(user.FriendId)}).
		Update("message_id", new_message.Id)

	return helper.SuccessResponse(c, 0, "Message Sent!", echo.Map{
		"message": request.Message,
		"id":      new_message.Id,
	})
}

func GetMessage(c echo.Context) error {
	claims, err := helper.JWT(c)

	if err != nil {
		return helper.ErrorResponse(c, 2, err.Error(), nil)
	}

	var user models.User

	if err := helper.Database.Db.Preload("Message", "created_at > ?", time.Now().Add(-(3*time.Hour))).First(&user, "token =?", claims.Token).Error; err != nil {
		return helper.ErrorResponse(c, 2, "No Messages", nil)
	}

	if user.Message.Id == 0 {
		return helper.ErrorResponse(c, 1, "No Messages", nil)
	}

	return helper.SuccessResponse(c, 0, "Your Message!", echo.Map{
		"message": echo.Map{
			"id":              user.Message.Id,
			"message":         user.Message.Message,
			"type_of_message": user.Message.TypeOfMessage,
			"created_time":    user.Message.CreatedAt,
		},
	})
}

func ChangeFriend(c echo.Context) error {
	claims, err := helper.JWT(c)

	if err != nil {
		return helper.ErrorResponse(c, 2, err.Error(), nil)
	}

	request := changeFriendStruct{}
	if err := c.Bind(&request); err != nil {
		return helper.ErrorResponse(c, 1, "Error Validation", nil)
	}

	if request.ChangeToUser == 0 {
		return helper.ErrorResponse(c, 1, "Friend number is not correct", nil)
	}

	var user models.User

	if err := helper.Database.Db.First(&user, "token =?", claims.Token).Error; err != nil {
		return helper.ErrorResponse(c, 2, "Something Went Wrong!", nil)
	}

	if request.Applying {
		if request.Accepting {
			if user.FriendId != 0 {
				helper.Database.Db.Model(&models.User{}).Where("id = ?", user.FriendId).Update("friend_id", 0)
			}
			var notification models.Notification

			helper.Database.Db.Where("id = ? AND extra_value = ?", request.ChangeToUser, user.Id).Find(&notification)

			if notification.Id == 0 {
				return helper.ErrorResponse(c, 1, "Something Went Wrong!", nil)
			}
			send_fcm_notification(notification.ExtraValue, "Friend Request", fmt.Sprintf("%d sent you a friend request.", user.PhoneNumber))

			err := apply_notification(int(user.Id), notification.Message, notification.UserId)
			if err != nil {
				return helper.ErrorResponse(c, 2, err.Error(), nil)
			}
		}
		helper.Database.Db.Where("id = ? AND extra_value = ?", request.ChangeToUser, user.Id).Delete(&models.Notification{})

		return helper.SuccessResponse(c, 0, "Success!", nil)
	}

	if user.FriendId != 0 {
		helper.Database.Db.Model(&models.User{}).Where("id = ?", user.FriendId).Update("friend_id", 0)
		helper.Database.Db.Model(&models.User{}).Where("id = ?", user.Id).Update("friend_id", 0)
	}

	var new_frined models.User

	helper.Database.Db.Select("id").First(&new_frined, "phone_number = ?", request.ChangeToUser)

	helper.Database.Db.Save(&models.Notification{
		Message:    "friend_request",
		UserId:     user.Id,
		ExtraValue: strconv.Itoa(int(new_frined.Id)),
	})
	return helper.SuccessResponse(c, 0, "Request Sent!", nil)
}

func Notification(c echo.Context) error {

	claims, err := helper.JWT(c)

	if err != nil {
		return helper.ErrorResponse(c, 2, err.Error(), nil)
	}

	var user models.User

	if err := helper.Database.Db.First(&user, "token =?", claims.Token).Error; err != nil {
		return helper.ErrorResponse(c, 2, "Something Went Wrong!", nil)
	}

	var notifications []models.Notification

	helper.Database.Db.Find(&notifications, "extra_value =?", user.Id)

	return helper.SuccessResponse(c, 0, "Notifications", echo.Map{
		"notifications": notifications,
	})
}

func send_fcm_notification(user_id interface{}, title, body string) error {
	fcm := []models.Fcm{}
	helper.Database.Db.Select("token").Where("user_id = ?", user_id).Find(&fcm)
	fcms := []string{}

	for _, v := range fcm {
		fcms = append(fcms, v.Token)
	}
	helper.Notification(fcms, title, body)
	return nil
}

func apply_notification(user_id int, method, extra_value interface{}) error {
	switch method {
	case "friend_request":
		helper.Database.Db.Model(&models.User{}).Where("id = ?", user_id).Update("friend_id", extra_value)
		helper.Database.Db.Model(&models.User{}).Where("id = ?", extra_value).Update("friend_id", user_id)
		return nil
	}
	return fmt.Errorf("not impleminated")
}

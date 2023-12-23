package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"gorm.io/gorm"

	"github.com/golang-jwt/jwt/v5"
)

var (
	Database  DbInstance
	serverKey = EnvVariable("FCM_TOKEN")
)

const (
	charset     = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	fcmEndpoint = "https://fcm.googleapis.com/fcm/send"
)

type DbInstance struct {
	Db *gorm.DB
}

type UserJwtClaims struct {
	Token string `json:"token"`
	Role  int    `json:"role"`
	jwt.RegisteredClaims
}

func EnvVariable(key string) string {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()

	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}
	value, ok := viper.Get(key).(string)

	if !ok {
		log.Fatalf("Invalid type assertion")
	}

	return value
}

func RandomString(n int) string {
	result := make([]byte, n)

	for i := range result {
		result[i] = charset[rand.Intn(len(charset)-1)]
	}

	return string(result)
}

func JWTAuthUser(token interface{}, user interface{}) error {
	switch token.(type) {
	case string:
		if err := Database.Db.First(&user, "token =?", token).Error; err != nil {
			return err
		}
	case uint:
		if err := Database.Db.First(&user, "id =?", token).Error; err != nil {
			return err
		}
	}
	return nil
}

func ErrorResponse(c echo.Context, result_type int, message interface{}, data echo.Map) error {
	return c.JSON(http.StatusBadRequest, echo.Map{"message": message, "result": result_type, "data": data})
}

func SuccessResponse(c echo.Context, result_type int, message interface{}, data echo.Map) error {
	return c.JSON(http.StatusOK, echo.Map{"message": message, "result": result_type, "data": data})
}

func JWT(c echo.Context) (UserJwtClaims, error) {
	authHeader := c.Request().Header.Get("Authorization")

	if authHeader == "" {
		return UserJwtClaims{}, fmt.Errorf("missing token")
	}

	if len(authHeader) < 10 {
		return UserJwtClaims{}, fmt.Errorf("token is not valid")
	}

	tokenString := authHeader[len("Bearer "):]

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(EnvVariable("SERECT")), nil
	})

	if err != nil || !token.Valid {
		return UserJwtClaims{}, fmt.Errorf("token is not valid")
	}

	claims := token.Claims.(jwt.MapClaims)

	return UserJwtClaims{
		Token: claims["token"].(string),
		Role:  int(claims["role"].(float64)),
	}, nil
}

func Notification(registrationIDs []string, title, body string) {
	payload := map[string]interface{}{
		"registration_ids": registrationIDs,
		"notification": map[string]interface{}{
			"title": title,
			"body":  body,
		},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	request, err := http.NewRequest("POST", fcmEndpoint, bytes.NewBuffer(payloadBytes))
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "key="+serverKey)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Error making HTTP request:", err)
		return
	}
	defer response.Body.Close()

	// bodyBytes, err := io.ReadAll(response.Body)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// bodyString := string(bodyBytes)
	// fmt.Println(bodyString)

	// fmt.Println(response.Body)
	// if response.StatusCode == http.StatusOK {
	// 	fmt.Println("Notification sent successfully.")
	// } else {
	// 	fmt.Println("Error sending notification. Status code:", response.StatusCode)
	// }
}

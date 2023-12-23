package routes

import (
	"khanakhaneho/controller"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Setuprouter(e *echo.Echo) {
	e.GET("/", home)

	e.POST("/login", controller.Login)

	e.POST("/me", controller.Me)
	e.POST("/change_friend", controller.ChangeFriend)
	e.POST("/notification", controller.Notification)
	e.POST("/send/message", controller.SendMessage)
	e.POST("/get/message", controller.GetMessage)
}

func home(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, Home!")
}

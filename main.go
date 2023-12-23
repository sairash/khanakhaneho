package main

import (
	"khanakhaneho/database"
	"khanakhaneho/routes"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.CORS())

	database.ConnectDb()

	routes.Setuprouter(e)

	e.Logger.Fatal(e.Start(":8080"))
}

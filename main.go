package main

import (
	"log/slog"
	"todo-api/config"
	"todo-api/routes"

	"github.com/labstack/echo/v4"
)

func main() {
	slog.Info("starting server")

	config.InitDB()

	e := echo.New()
	routes.RegisterRoutes(e)

	e.Logger.Fatal(e.Start(":8080"))
}

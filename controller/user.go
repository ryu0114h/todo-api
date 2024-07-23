package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type UserController interface {
	CreateUser(ctx echo.Context) error
}

type userController struct {
}

func NewUserController() UserController {
	return &userController{}
}

func (c *userController) CreateUser(ctx echo.Context) error {

	return ctx.JSON(http.StatusCreated, nil)
}

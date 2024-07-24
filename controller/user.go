package controller

import (
	"net/http"
	"todo-api/controller/request"
	"todo-api/controller/response"
	"todo-api/usecase"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type UserController interface {
	CreateUser(ctx echo.Context) error
}

type userController struct {
	validate    *validator.Validate
	userUseCase usecase.UserUseCase
}

func NewUserController(validate *validator.Validate, userUseCase usecase.UserUseCase) UserController {
	return &userController{
		validate:    validate,
		userUseCase: userUseCase,
	}
}

func (c *userController) CreateUser(ctx echo.Context) error {
	var requestBody request.CreateUserRequestBody
	if err := ctx.Bind(&requestBody); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	// Validation
	if err := c.validate.Struct(requestBody); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	createdUser, err := c.userUseCase.CreateUser(
		requestBody.Username,
		requestBody.Email,
		requestBody.Password,
		requestBody.Role,
		requestBody.CompanyIds,
	)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "could not create user"})
	}

	return ctx.JSON(http.StatusCreated, response.NewCreateUserResponseBody(createdUser))
}

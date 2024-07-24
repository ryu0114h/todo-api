package controller

import (
	"net/http"
	"os"
	"time"
	"todo-api/controller/request"
	"todo-api/controller/response"
	"todo-api/usecase"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type AuthController interface {
	Login(ctx echo.Context) error
}

type authController struct {
	validate    *validator.Validate
	authUseCase usecase.AuthUseCase
}

func NewAuthController(validate *validator.Validate, authUseCase usecase.AuthUseCase) AuthController {
	return &authController{
		validate:    validate,
		authUseCase: authUseCase,
	}
}

func (c *authController) Login(ctx echo.Context) error {
	var requestBody request.LoginRequestBody
	if err := ctx.Bind(&requestBody); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	// Validation
	if err := c.validate.Struct(requestBody); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	loginUser, err := c.authUseCase.Login(
		requestBody.Username,
		requestBody.Password,
	)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "could not create auth"})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  loginUser.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "JWT Secret not found"})
	}

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Error while generating token"})
	}

	return ctx.JSON(http.StatusOK, response.NewLoginResponseBody(tokenString))
}

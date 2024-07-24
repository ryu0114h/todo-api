package middleware

import (
	"net/http"
	"todo-api/model"

	"github.com/labstack/echo/v4"
)

// Adminチェックミドルウェア
func AdminAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			user := ctx.Get("user").(*model.User)
			if user.Role != "admin" {
				return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "not authorization"})
			}

			return next(ctx)
		}
	}
}

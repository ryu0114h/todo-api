package middleware

import (
	"net/http"
	"os"
	"strings"
	"todo-api/repository"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

// 認証ミドルウェア
func Auth(userRepository repository.UserRepository) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			tokenString := ctx.Request().Header.Get("Authorization")
			if tokenString == "" {
				return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "Missing token"})
			}

			// "Bearer " プレフィックスを削除
			tokenString = strings.TrimPrefix(tokenString, "Bearer ")

			// JWT シークレットの取得
			jwtSecret := os.Getenv("JWT_SECRET")
			if jwtSecret == "" {
				return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "JWT Secret not found"})
			}

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
				}
				return []byte(jwtSecret), nil
			})

			if err != nil || !token.Valid {
				return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok || !token.Valid {
				return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token claims"})
			}
			id := claims["id"].(float64)
			user, err := userRepository.GetUser(uint(id))
			if err != nil {
				return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve user"})
			}
			ctx.Set("user", user)

			return next(ctx)
		}
	}
}

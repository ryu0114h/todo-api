package middleware

import (
	"errors"
	"net/http"
	"strconv"
	myErrors "todo-api/errors"
	"todo-api/model"
	"todo-api/repository"

	"github.com/labstack/echo/v4"
)

// Companyに関する認証ミドルウェア
// ログインしているユーザがCompanyに所属していない場合、not found
func CompanyAuth(companyUserRepository repository.CompanyUserRepository) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			companyId, err := strconv.ParseUint(ctx.Param("company_id"), 10, 64)
			if err != nil {
				return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "company_id is bad request"})
			}
			user := ctx.Get("user").(*model.User)

			_, err = companyUserRepository.GetCompanyUser(uint(companyId), uint(user.ID))
			if err != nil {
				if errors.Is(err, myErrors.ErrNotFound) {
					return ctx.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
				}

				return ctx.JSON(http.StatusBadRequest, nil)
			}

			return next(ctx)
		}
	}
}

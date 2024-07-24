package middleware

import (
	"fmt"
	"log/slog"

	"github.com/labstack/echo/v4"
)

// リクエストのメソッドとURIをログ出力するミドルウェア
// echo.MiddlewareFunc を返す
func Logging() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			// ログを出力
			slog.Info(fmt.Sprintf("%s %s", ctx.Request().Method, ctx.Request().RequestURI))

			return next(ctx)
		}
	}
}

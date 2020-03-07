package middleware

import (
	"operator/config"

	"github.com/labstack/echo"
)

func DbMiddleware(db config.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("db", db.Db)
			return next(c)
		}
	}
}

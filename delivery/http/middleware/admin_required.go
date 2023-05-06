package middleware

import (
	deliveryHttp "github.com/behnambm/assignment/utils/http"
	"github.com/labstack/echo"
	"net/http"
)

func AdminRequired() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if !deliveryHttp.IsAdmin(c) {
				return c.JSON(http.StatusForbidden, echo.Map{"error": "access denied"})
			}
			return next(c)
		}
	}
}

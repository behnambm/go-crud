package middleware

import (
	"github.com/behnambm/assignment/constants"
	"github.com/labstack/echo"
	"net/http"
)

func LoginRequired() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			isAuthenticated, ok := c.Get(constants.IsAuthenticatedKey).(bool)
			if !ok {
				return c.JSON(http.StatusInternalServerError, echo.Map{"error": "could not authenticate"})
			}
			if !isAuthenticated {
				return c.JSON(http.StatusForbidden, echo.Map{"error": "access denied"})
			}
			return next(c)
		}
	}
}

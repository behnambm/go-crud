package http

import (
	"github.com/behnambm/go-crud/constants"
	"github.com/labstack/echo/v4"
)

func IsAuthenticated(c echo.Context) bool {
	authenticated, ok := c.Get(constants.IsAuthenticatedKey).(bool)
	if !ok {
		return false
	}
	return authenticated
}

func IsAdmin(c echo.Context) bool {
	isAdmin, ok := c.Get(constants.IsAdminKey).(bool)
	if !ok {
		return false
	}
	return isAdmin
}

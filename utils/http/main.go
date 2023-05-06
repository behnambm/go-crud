package http

import (
	"github.com/behnambm/assignment/constants"
	"github.com/labstack/echo"
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

package http

import (
	"github.com/behnambm/assignment/delivery/http/middleware"
	"github.com/labstack/echo"
)

func IsAuthenticated(c echo.Context) bool {
	authenticated, ok := c.Get(middleware.IsAuthenticatedKey).(bool)
	if !ok {
		return false
	}
	return authenticated
}

func IsAdmin(c echo.Context) bool {
	isAdmin, ok := c.Get(middleware.IsAdminKey).(bool)
	if !ok {
		return false
	}
	return isAdmin
}

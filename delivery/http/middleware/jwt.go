package middleware

import (
	"github.com/behnambm/assignment/service/auth"
	"github.com/behnambm/assignment/service/user"
	"github.com/labstack/echo"
	"net/http"
)

func Auth(userSrv user.Service, authSrv auth.Service) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusBadRequest, echo.Map{"error": "no auth header provided"})
			}

			claim, valid := authSrv.IsValidWithClaim(authHeader)
			if !valid {
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "access denied"})
			}
			uid, ok := claim["uid"].(float64)
			if !ok {
				return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid token"})
			}
			currentUser, userErr := userSrv.GetUserFromID(int(uid))

			if userErr != nil {
				return c.JSON(http.StatusNotFound, echo.Map{"error": "user not found"})
			}
			c.Set("current_user", currentUser)
			return next(c)
		}
	}
}

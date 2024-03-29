package middleware

import (
	"github.com/behnambm/go-crud/constants"
	"github.com/behnambm/go-crud/service/auth"
	"github.com/behnambm/go-crud/service/user"
	"github.com/labstack/echo/v4"
)

func Auth(userSrv user.Service, authSrv auth.Service) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			isAuthenticated := false

			authHeader := c.Request().Header.Get("Authorization")
			if authHeader != "" {
				claim, valid := authSrv.IsValidWithClaim(authHeader)
				if valid {
					uid, ok := claim["uid"].(float64)
					if ok {
						currentUser, userErr := userSrv.GetUserFromID(int(uid))
						if userErr == nil {
							isAuthenticated = true
							c.Set(constants.CurrentUserKey, currentUser)
							c.Set(constants.IsAdminKey, currentUser.IsAdmin)
						}
					}
				}
			}
			c.Set(constants.IsAuthenticatedKey, isAuthenticated)

			return next(c)
		}
	}
}

package middleware

import (
	"github.com/behnambm/assignment/service/auth"
	"github.com/behnambm/assignment/service/user"
	"github.com/labstack/echo"
)

const (
	IsAuthenticatedKey = "IsAuthenticated"
	CurrentUserKey     = "CurrentUser"
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
							c.Set(CurrentUserKey, currentUser)
						}
					}
				}
			}
			c.Set(IsAuthenticatedKey, isAuthenticated)

			return next(c)
		}
	}
}

package middlewares

import (
	"budget-app/common"
	"budget-app/internal/model"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"strings"
)

type AppMiddleware struct {
	Logger echo.Logger
	Db     *gorm.DB
}

func (appMiddleware *AppMiddleware) AuthenticateMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Add("Vary", "Authorization")
		authHeader := c.Request().Header.Get("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") == false {
			return common.SendUnauthorizedResponse(c, "Please provide bearer token")
		}
		accessToken := strings.Split(authHeader, " ")[1]
		fmt.Println("Here is auth header: ", accessToken)
		claims, err := common.ParseJwtSignedAccessToken(accessToken)
		if err != nil {
			return common.SendUnauthorizedResponse(c, "Unauthorized user")
		}

		if common.IsExpired(claims) {
			return common.SendUnauthorizedResponse(c, "Token expired")
		}
		var user model.User
		result := appMiddleware.Db.First(&user, claims.ID)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return common.SendUnauthorizedResponse(c, "Invalid token")
		}
		if result.Error != nil {
			return common.SendUnauthorizedResponse(c, "Unauthorized user")
		}
		c.Set("user", user)
		return next(c)
	}
}

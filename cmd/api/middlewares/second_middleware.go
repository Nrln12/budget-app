package middlewares

import (
	"fmt"
	"github.com/labstack/echo/v4"
)

func SecondMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	fmt.Printf("we are in second-middleware\n")
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderServer, "Echo/3.0")
		return next(c)
	}
}

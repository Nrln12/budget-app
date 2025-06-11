package handlers

import (
	"budget-app/internal/mailer"
	"errors"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Handler struct {
	Database *gorm.DB
	Logger   echo.Logger
	Mailer   mailer.Mailer
}

func (h *Handler) BindRequestBody(c echo.Context, payload interface{}) error {
	if err := (&echo.DefaultBinder{}).BindBody(c, payload); err != nil {
		c.Logger().Error(err)
		return errors.New("failed to bind request body: " + err.Error())
	}
	return nil
}

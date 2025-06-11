package common

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type JsonErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func SendErrorResponse(c echo.Context, message string, statusCode int) error {
	return c.JSON(statusCode, JsonErrorResponse{
		Success: false,
		Message: message,
	})
}

func SendSuccessResponse(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, data)
}

func SendBadRequestResponse(c echo.Context, message string) error {
	return SendErrorResponse(c, message, http.StatusBadRequest)
}

func SendNotFoundResponse(c echo.Context, message string) error {
	return SendErrorResponse(c, message, http.StatusNotFound)
}
func SendInternalServerErrorResponse(c echo.Context, message string) error {
	return SendErrorResponse(c, message, http.StatusInternalServerError)
}

func SendUnauthorizedResponse(c echo.Context, message string) error {
	return SendErrorResponse(c, message, http.StatusUnauthorized)
}

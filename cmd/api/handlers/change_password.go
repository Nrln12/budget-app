package handlers

import (
	"budget-app/cmd/api/request"
	"budget-app/cmd/api/services"
	"budget-app/common"
	"budget-app/internal/model"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *Handler) ChangePasswordHandler(c echo.Context) error {
	user, ok := c.Get("user").(model.User)
	if !ok {
		return common.SendInternalServerErrorResponse(c, "Authentication failed")
	}
	payload := new(request.ChangePasswordRequest)
	if err := h.BindRequestBody(c, payload); err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	validationErrors := h.ValidateRequestBody(c, *payload)
	if validationErrors != nil {
		return c.JSON(http.StatusBadRequest, validationErrors)
	}

	// check is current password correct
	if common.CheckPasswordHash(payload.CurrentPassword, user.Password) == false {
		return common.SendBadRequestResponse(c, "Incorrect Password")
	}
	userService := services.NewUserService(h.Database)
	err := userService.ChangePassword(payload.NewPassword, user)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, err.Error())
	}
	return common.SendSuccessResponse(c, "Password changed successfully")
}

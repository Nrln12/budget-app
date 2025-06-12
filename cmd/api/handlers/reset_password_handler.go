package handlers

import (
	"budget-app/cmd/api/request"
	"budget-app/cmd/api/services"
	"budget-app/common"
	"budget-app/internal/mailer"
	"encoding/base64"
	"errors"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
	"net/url"
	"os"
)

func (h *Handler) ForgotPasswordHandler(c echo.Context) error {
	// bind request body
	payload := new(request.ForgotPasswordRequest)
	if err := h.BindRequestBody(c, payload); err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}
	validationErrors := h.ValidateRequestBody(c, payload)
	if validationErrors != nil {
		return c.JSON(http.StatusBadRequest, validationErrors)
	}
	userService := services.NewUserService(h.Database)
	appTokenService := services.NewAppTokenService(h.Database)

	// check user
	user, err := userService.GetUserByEmail(payload.Email)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return common.SendNotFoundResponse(c, "User not found")
	}
	if err != nil {
		return common.SendInternalServerErrorResponse(c, "Unexpected error happened")
	}

	// generate token
	token, err := appTokenService.GenerateResetPasswordToken(*user)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, "Unexpected error happened")
	}

	encodedEmail := base64.RawURLEncoding.EncodeToString([]byte(user.Email))
	frontEndUrl, err := url.Parse(payload.FrontendUrl)
	if err != nil {
		return common.SendBadRequestResponse(c, "Invalid frontend url")
	}
	query := url.Values{}
	query.Set("email", encodedEmail)
	query.Set("token", token.Token)
	frontEndUrl.RawQuery = query.Encode()

	mailData := mailer.EmailData{
		Subject: "Welcome To " + os.Getenv("APP_NAME"),
		Meta: struct {
			Token       string
			FrontendUrl string
		}{
			Token:       token.Token,
			FrontendUrl: frontEndUrl.String(),
		},
	}
	err = h.Mailer.Send(payload.Email, "forgot_password.html", mailData)
	if err != nil {
		h.Logger.Error(err)
	}
	return common.SendSuccessResponse(c, "Forgot password email sent successfully")
}

func (h *Handler) ResetPasswordHandler(c echo.Context) error {
	// bind request body
	payload := new(request.ResetPasswordRequest)
	if err := h.BindRequestBody(c, payload); err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	validationErrors := h.ValidateRequestBody(c, payload)
	if validationErrors != nil {
		return c.JSON(http.StatusBadRequest, validationErrors)
	}
	// to encode url format
	email, err := base64.RawURLEncoding.DecodeString(payload.Meta)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, "Unexpected error happened")
	}

	userService := services.NewUserService(h.Database)
	appTokenService := services.NewAppTokenService(h.Database)

	user, err := userService.GetUserByEmail(string(email))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return common.SendNotFoundResponse(c, "Invalid password reset token")
	}
	if err != nil {
		return common.SendInternalServerErrorResponse(c, "Unexpected error happened")
	}

	token, err := appTokenService.ValidateResetPasswordToken(*user, payload.Token)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, err.Error())
	}

	err = userService.ChangePassword(payload.Password, *user)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, err.Error())
	}

	appTokenService.InvalidateToken(user.ID, *token)
	return common.SendSuccessResponse(c, "Reset password successfully")
}

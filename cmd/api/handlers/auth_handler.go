package handlers

import (
	"budget-app/cmd/api/request"
	"budget-app/cmd/api/services"
	"budget-app/common"
	"budget-app/internal/mailer"
	"budget-app/internal/model"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
)

func (h *Handler) RegisterHandler(c echo.Context) error {
	payload := new(request.RegisterUserRequest)
	if err := h.BindRequestBody(c, payload); err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}
	c.Logger().Info("Payload: ", payload)

	validationErrors := h.ValidateRequestBody(c, *payload)
	if validationErrors != nil {
		return c.JSON(http.StatusBadRequest, validationErrors)
	}

	userService := services.NewUserService(h.Database)
	_, err := userService.GetUserByEmail(payload.Email)

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		h.Logger.Error("DB error: ", err)
		return c.JSON(http.StatusBadRequest, "user already exists")
	}

	user, err := userService.CreateUser(*payload)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "user creation failed")
	}

	mailData := mailer.EmailData{
		Subject: "Hello from Budget App",
		Meta:    nil,
	}

	err = h.Mailer.Send(payload.Email, "hello.html", mailData)
	if err != nil {
		h.Logger.Error("Error sending email: ", err)
	}

	fmt.Printf("Created user: %v\n", user)
	return c.JSON(http.StatusOK, "ok")
}

func (h *Handler) LoginHandler(c echo.Context) error {
	userService := services.NewUserService(h.Database)
	// bind our data
	payload := new(request.LoginRequest)
	if err := h.BindRequestBody(c, payload); err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	// validate the data
	validationErrors := h.ValidateRequestBody(c, *payload)
	if validationErrors != nil {
		return c.JSON(http.StatusBadRequest, validationErrors)
	}

	// check email exists
	user, err := userService.GetUserByEmail(payload.Email)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.JSON(http.StatusBadRequest, "Invalid credentials")
	}
	fmt.Printf("Found user: %v\n", *user)

	// check password
	if common.CheckPasswordHash(payload.Password, user.Password) == false {
		return c.JSON(http.StatusBadRequest, "Invalid credentials")
	}
	// sending response with user token
	accessToken, refreshToken, err := common.GenerateJwt(*user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "internal server error")
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
		"user": map[string]interface{}{
			"id":    user.ID,
			"email": user.Email,
		},
	})
}

func (h *Handler) GetAuthenticatedUserHandler(c echo.Context) error {
	user, ok := c.Get("user").(model.User)
	if !ok {
		return common.SendInternalServerErrorResponse(c, "Authentication failed")
	}
	return common.SendSuccessResponse(c, user)
}

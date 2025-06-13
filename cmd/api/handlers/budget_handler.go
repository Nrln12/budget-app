package handlers

import (
	"budget-app/cmd/api/request"
	"budget-app/cmd/api/services"
	"budget-app/common"
	"budget-app/internal/model"
	"github.com/labstack/echo/v4"
)

func (h *Handler) CreateBudget(c echo.Context) error {
	user, ok := c.Get("user").(model.User)
	if !ok {
		return common.SendUnauthorizedResponse(c, "User is unauthorized")
	}
	payload := new(request.CreateBudgetRequest)
	if err := h.BindRequestBody(c, payload); err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	validationErrors := h.ValidateRequestBody(c, *payload)
	if validationErrors != nil {
		return common.SendBadRequestResponse(c, "Validation failed")
	}

	budgetService := services.NewBudgetService(h.Database)
	//categoryService := services.NewCategoryService(h.Database)

	_, err := budgetService.CreateBudget(payload, user.ID)
	if err != nil {
		c.Logger().Error(err)
		return common.SendInternalServerErrorResponse(c, "Budget could not be created")
	}

	//categoryService
	return common.SendSuccessResponse(c, "Budget created")
}

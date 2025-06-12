package handlers

import (
	"budget-app/cmd/api/request"
	"budget-app/cmd/api/services"
	"budget-app/common"
	"budget-app/internal/app_errors"
	"budget-app/internal/model"
	"errors"
	"github.com/labstack/echo/v4"
)

func (h *Handler) GetCategories(c echo.Context) error {
	var categories []*model.Category
	categoryService := services.NewCategoryService(h.Database)
	pagination := common.NewPageResponse(categories, c.Request(), h.Database)
	categoryPage, err := categoryService.GetCategories(categories, pagination)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, err.Error())
	}
	return common.SendSuccessResponse(c, categoryPage)
}

func (h *Handler) CreateCategory(c echo.Context) error {
	_, ok := c.Get("user").(model.User)
	if !ok {
		return common.SendInternalServerErrorResponse(c, "User is not logged in")
	}
	// bind request body
	payload := new(request.CreateCategoryRequest)
	if err := h.BindRequestBody(c, payload); err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	// validation
	validationErrors := h.ValidateRequestBody(c, *payload)
	if validationErrors != nil {
		return common.SendBadRequestResponse(c, "Validation Failed")
	}

	categoryService := services.NewCategoryService(h.Database)
	category, err := categoryService.CreateCategory(payload)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, err.Error())
	}
	return common.SendSuccessResponse(c, category)
}

func (h *Handler) DeleteCategory(c echo.Context) error {
	_, ok := c.Get("user").(model.User)
	if !ok {
		return common.SendInternalServerErrorResponse(c, "User is not logged in")
	}
	var categoryId request.IdParamRequest
	err := (&echo.DefaultBinder{}).BindPathParams(c, &categoryId)
	if err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}
	categoryService := services.NewCategoryService(h.Database)
	err = categoryService.DeleteById(categoryId.Id)
	if err != nil {
		if errors.Is(err, app_errors.NewNotFoundError(err.Error())) {
			return common.SendNotFoundResponse(c, err.Error())
		}
		return common.SendInternalServerErrorResponse(c, err.Error())
	}
	return common.SendSuccessResponse(c, nil)
}

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

func (h *Handler) GetBudget(c echo.Context) error {
	user, _ := c.Get("user").(model.User)
	var budgets []*model.Budget
	budgetService := services.NewBudgetService(h.Database)
	query := h.Database.Preload("Categories").Scopes(common.WhereUserIdScope(user.ID))
	paginator := common.NewPageResponse(budgets, c.Request(), query)
	budgetPage, err := budgetService.GetBudget(query, budgets, paginator)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, err.Error())
	}
	return common.SendSuccessResponse(c, &budgetPage)
}

func (h *Handler) CreateBudget(c echo.Context) error {
	user, _ := c.Get("user").(model.User)
	payload := new(request.CreateBudgetRequest)
	if err := h.BindRequestBody(c, payload); err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	validationErrors := h.ValidateRequestBody(c, *payload)
	if validationErrors != nil {
		return common.SendBadRequestResponse(c, "Validation failed")
	}

	budgetService := services.NewBudgetService(h.Database)
	categoryService := services.NewCategoryService(h.Database)

	budget, err := budgetService.CreateBudget(payload, user.ID)
	if err != nil {
		c.Logger().Error(err)
		return common.SendInternalServerErrorResponse(c, "Budget could not be created")
	}

	categories, err := categoryService.GetMultipleCategories(payload.Categories)
	if err != nil {
		c.Logger().Error(err)
		return common.SendInternalServerErrorResponse(c, "Budget could not be created")
	}

	err = budgetService.Db.Model(&budget).Association("Categories").Replace(categories)
	if err != nil {
		c.Logger().Error(err)
		return common.SendInternalServerErrorResponse(c, "Budget could not be created")
	}
	budget.Categories = categories
	return common.SendSuccessResponse(c, budget)
}

func (h *Handler) UpdateBudget(c echo.Context) error {
	user, _ := c.Get("user").(model.User)

	var budgetId request.IdParamRequest
	err := (&echo.DefaultBinder{}).BindPathParams(c, &budgetId)
	if err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	budgetService := services.NewBudgetService(h.Database)
	categoryService := services.NewCategoryService(h.Database)

	budget, err := budgetService.GetBudgetById(budgetId.Id)
	if err != nil {
		if errors.Is(err, app_errors.NewNotFoundError(err.Error())) {
			return common.SendNotFoundResponse(c, err.Error())
		}
		return common.SendInternalServerErrorResponse(c, err.Error())
	}

	if user.ID != budget.UserId {
		return common.SendNotFoundResponse(c, "Budget not found")
	}

	payload := new(request.UpdateBudgetRequest)
	if err := h.BindRequestBody(c, payload); err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	validationErrors := h.ValidateRequestBody(c, *payload)
	if validationErrors != nil {
		return common.SendBadRequestResponse(c, "Validation failed")
	}
	updateBudget, err := budgetService.UpdateBudget(budget, payload, user.ID)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, err.Error())
	}

	if payload.Categories != nil {
		categories, _ := categoryService.GetMultipleCategories(payload.Categories)
		err = budgetService.Db.Model(&updateBudget).Association("Categories").Replace(categories)
		if err != nil {
			c.Logger().Error(err)
			return common.SendInternalServerErrorResponse(c, "Budget could not be updated")
		}
		updateBudget.Categories = categories
	}
	return common.SendSuccessResponse(c, updateBudget)
}

func (h *Handler) DeleteBudget(c echo.Context) error {
	user, _ := c.Get("user").(model.User)

	var budgetId request.IdParamRequest
	err := (&echo.DefaultBinder{}).BindPathParams(c, &budgetId)
	if err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	budgetService := services.NewBudgetService(h.Database)
	budget, _ := budgetService.GetBudgetById(budgetId.Id)
	if budget == nil {
		return common.SendNotFoundResponse(c, "Budget not found")
	}

	if user.ID != budget.UserId {
		return common.SendNotFoundResponse(c, "Budget not found")
	}

	query := h.Database.Scopes(common.WhereUserIdScope(user.ID))
	query.Delete(model.Budget{}, budget.ID)
	return common.SendSuccessResponse(c, "Budget deleted successfully")
}

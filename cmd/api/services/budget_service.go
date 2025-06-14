package services

import (
	"budget-app/cmd/api/request"
	"budget-app/common"
	"budget-app/internal/app_errors"
	"budget-app/internal/model"
	"errors"
	"fmt"
	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
	"strings"
	"time"
)

type BudgetService struct {
	Db *gorm.DB
}

func NewBudgetService(db *gorm.DB) *BudgetService {
	return &BudgetService{Db: db}
}

func (budgetService *BudgetService) CreateBudget(budgetRequest *request.CreateBudgetRequest, userId uint) (*model.Budget, error) {
	slug := strings.ToLower(budgetRequest.Title)
	slug = strings.Replace(slug, " ", "_", -1)
	budget := &model.Budget{
		Amount:      budgetRequest.Amount,
		UserId:      userId,
		Title:       budgetRequest.Title,
		Slug:        slug,
		Description: budgetRequest.Description,
	}
	if budgetRequest.Date == "" {
		budget.Date = time.Now()
	}
	budget.Month = uint(budget.Date.Month())
	budget.Year = uint16(budget.Date.Year())
	budgetExists, err := budgetService.budgetExists(userId, budget.Month, budget.Year, budget.Slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result := budgetService.Db.Create(budget)
			if result.Error != nil {
				return nil, result.Error
			}
			return budget, nil
		}
		return nil, err
	}
	return budgetExists, nil
}

func (budgetService *BudgetService) budgetExists(userId uint, month uint, year uint16, slug string) (*model.Budget, error) {
	budget := model.Budget{}
	result := budgetService.Db.
		Where("user_id = ? AND month = ? AND year = ? AND slug = ?",
			userId, month, year, slug).First(&budget)
	if result.Error != nil {
		return nil, result.Error
	}
	return &budget, nil
}

func (budgetService *BudgetService) countForYearAndMonthAndSlugAndUserIdExcludeBudgetId(userId uint, month uint, year uint16, slug string, budgetId uint) int64 {
	var count int64
	budgetService.Db.Model(&model.Budget{}).
		Where("budgets.user_id = ? AND budgets.month = ? AND budgets.year = ? AND budgets.slug = ? AND budgets.id != ?",
			userId, month, year, slug, budgetId).
		Count(&count)
	return count
}

func (budgetService *BudgetService) GetBudget(query *gorm.DB, budget []*model.Budget, pagination *common.PageResponse) (*common.PageResponse, error) {
	query.Scopes(pagination.Paginate()).Find(&budget)
	pagination.Items = budget
	return pagination, nil
}

func (budgetService *BudgetService) GetBudgetById(id uint) (*model.Budget, error) {
	var budget model.Budget
	result := budgetService.Db.First(&budget, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, app_errors.NewNotFoundError("Budget not found")
		}
		return nil, errors.New("failed to get budget")
	}
	return &budget, nil
}

func (budgetService *BudgetService) UpdateBudget(budget *model.Budget, payload *request.UpdateBudgetRequest, userId uint) (*model.Budget, error) {
	if payload.Date != "" {
		timeParsed, err := time.Parse(time.DateOnly, payload.Date)
		if err != nil {
			return nil, err
		}
		budget.Date = timeParsed
	}
	if payload.Amount > 0 {
		budget.Amount = payload.Amount
	}

	if payload.Description != nil {
		budget.Description = payload.Description
	}

	if payload.Title != "" {
		budget.Title = payload.Title
		slug := strings.ToLower(payload.Title)
		slug = strings.Replace(slug, " ", "_", -1)
		budget.Slug = slug
	}

	count := budgetService.countForYearAndMonthAndSlugAndUserIdExcludeBudgetId(userId, budget.Month, budget.Year, budget.Slug, budget.ID)
	log.Info("Passed parameters: ", userId, budget.Month, budget.Year, budget.Slug, budget.ID)
	fmt.Println("count is: ", count)
	if count > 0 {
		return nil, errors.New("budget with these details already exists")
	}
	budgetService.Db.Updates(budget)
	return budget, nil
}

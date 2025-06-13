package services

import (
	"budget-app/cmd/api/request"
	"budget-app/internal/model"
	"errors"
	"gorm.io/gorm"
	"strings"
	"time"
)

type BudgetService struct {
	db *gorm.DB
}

func NewBudgetService(db *gorm.DB) *BudgetService {
	return &BudgetService{db: db}
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
			result := budgetService.db.Create(budget)
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
	result := budgetService.db.
		Where("user_id = ? AND month = ? AND year = ? AND slug = ?",
			userId, month, year, slug).First(&budget)
	if result.Error != nil {
		return nil, result.Error
	}
	return &budget, nil
}

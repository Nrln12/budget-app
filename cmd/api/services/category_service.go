package services

import (
	"budget-app/cmd/api/request"
	"budget-app/common"
	"budget-app/internal/app_errors"
	"budget-app/internal/model"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
)

type CategoryService struct {
	Db *gorm.DB
}

func NewCategoryService(db *gorm.DB) *CategoryService {
	return &CategoryService{Db: db}
}

func (categoryService *CategoryService) GetCategories(categories []*model.Category, pagination *common.PageResponse) (*common.PageResponse, error) {
	//result := categoryService.Db.Find(&categories)
	result := categoryService.Db.Scopes(pagination.Paginate()).Find(&categories)
	pagination.Items = categories
	if result.Error != nil {
		return nil, errors.New("failed to get the categories")
	}
	return pagination, nil
}

func (categoryService *CategoryService) CreateCategory(categoryRequest *request.CreateCategoryRequest) (*model.Category, error) {
	slug := strings.ToLower(categoryRequest.Name)
	slug = strings.Replace(slug, " ", "_", -1)
	category := &model.Category{
		Slug:     slug,
		Name:     categoryRequest.Name,
		IsCustom: categoryRequest.IsCustom,
	}
	result := categoryService.Db.Where(model.Category{Slug: slug, Name: category.Name}).Assign(model.Category{IsCustom: categoryRequest.IsCustom}).FirstOrCreate(category)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return nil, errors.New("category already exists")
		}
		return nil, errors.New("failed to create the category")
	}
	return category, nil
}

func (categoryService *CategoryService) GetById(id uint) (*model.Category, error) {
	var category *model.Category
	result := categoryService.Db.First(&category, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, app_errors.NewNotFoundError("Category not found")
		}
		return nil, errors.New("failed to get the category")
	}
	return category, nil
}

func (categoryService *CategoryService) DeleteById(id uint) error {
	var category *model.Category
	category, err := categoryService.GetById(id)
	if err != nil {
		return err
	}
	categoryService.Db.Delete(category)
	return nil
}

func (c CategoryService) AssociateUserToCategories(user *model.User, categories []*model.Category) error {
	if user != nil && categories != nil && len(categories) > 0 {
		var userCategories []*model.UserCategory
		for _, category := range categories {
			userCategories = append(userCategories, &model.UserCategory{
				UserId:     user.ID,
				CategoryId: category.ID,
			})
		}
		result := c.Db.Clauses(clause.OnConflict{DoNothing: true}).Create(userCategories)
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}

func (c CategoryService) GetMultipleCategories(categoryIds []uint) ([]*model.Category, error) {
	var categories []*model.Category
	result := c.Db.Where("id in ?", categoryIds).Find(&categories)
	if result.Error != nil {
		return nil, result.Error
	}
	return categories, nil
}

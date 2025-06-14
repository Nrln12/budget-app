package request

type CreateCategoryRequest struct {
	Name     string `json:"name" validate:"required"`
	IsCustom bool   `json:"is_custom" default:"true"`
}

type UserCategoryRequest struct {
	Categories []uint `json:"categories" validate:"required"`
}

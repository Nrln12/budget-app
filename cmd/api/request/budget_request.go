package request

type CreateBudgetRequest struct {
	Categories  []uint64 `json:"categories" validate:"required,dive,min=1"`
	Amount      float64  `json:"amount" validate:"required,numeric,min=1"`
	Date        string   `json:"date,omitempty" validate:"omitempty,datetime=2006-01-02"`
	Title       string   `json:"title" validate:"required,min=2,max=255"`
	Description *string  `json:"description" validate:"omitempty,min=2,max=500"`
}

type UpdateBudgetRequest struct {
	Categories []uint64 `json:"categories" validate:"omitempty,dive,min=1"`
	Amount     float64  `json:"amount" validate:"omitempty,required,numeric,min=1"`
	Date       string   `json:"date,omitempty" validate:"omitempty,datetime=2006-01-02"`
	Title      string   `json:"title" validate:"omitempty,required,min=2,max=255"`
}

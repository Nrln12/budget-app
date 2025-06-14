package model

import "time"

type Budget struct {
	BaseModel
	Title       string      `json:"name" gorm:"index;type:varchar(255);not null"`
	Slug        string      `json:"slug" gorm:"index;type:varchar(255);not null"`
	Description *string     `json:"description" gorm:"type:text"`
	UserId      uint        `json:"user_id" gorm:"not null;column:user_id;unique_index:user_id_slug_year_month"`
	Amount      float64     `json:"amount" gorm:"type:decimal(10,2);not null"`
	Categories  []*Category `json:"categories" gorm:"constraint:OnDelete:CASCADE;many2many:budget_categories"`
	Date        time.Time   `json:"date" gorm:"type:datetime;not null"`
	Month       uint        `json:"month" gorm:"type:TINYINT;UNSIGNED;not null;index:idx_month_year;unique_index:user_id_slug_year_month"`
	Year        uint16      `json:"year" gorm:"type:INT;UNSIGNED;not null;index:idx_month_year;unique_index:user_id_slug_year_month"`
}

func (Budget) TableName() string {
	return "budgets"
}

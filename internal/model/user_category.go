package model

type UserCategory struct {
	UserId     uint `json:"user_id" gorm:"primary_key;column:user_id"`
	CategoryId uint `json:"category_id" gorm:"primary_key;column:category_id"`
}

func (UserCategory) TableName() string {
	return "user_categories"
}

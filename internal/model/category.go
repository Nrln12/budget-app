package model

type Category struct {
	BaseModel
	Name     string  `gorm:"unique;type:varchar(255);not null" json:"name"`
	Slug     string  `gorm:"unique;type:varchar(255);not null" json:"slug"`
	IsCustom bool    `gorm:"type:bool;default:false" json:"is_custom"`
	Users    []*User `gorm:"constraint:OnDelete:CASCADE;many2many:user_categories;joinForeignKey:CategoryId;joinReferences:UserId" json:"users,omitempty"`
}

func (Category) TableName() string {
	return "categories"
}

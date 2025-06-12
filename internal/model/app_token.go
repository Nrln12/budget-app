package model

import "time"

type AppToken struct {
	BaseModel
	Token     string    `json:"-" gorm:"index;type:varchar(255);not null"`
	TargetId  uint      `json:"target_id" gorm:"index;not null"`
	Type      string    `json:"-" gorm:"index;not null;type:varchar(255)"`
	Used      bool      `json:"-" gorm:"index;not null;type:boolean"`
	ExpiresAt time.Time `json:"-" gorm:"index;not null;type:datetime"`
}

func (AppToken) tableName() string {
	return "app_tokens"
}

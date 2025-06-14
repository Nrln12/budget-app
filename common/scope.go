package common

import (
	"fmt"
	"gorm.io/gorm"
)

func WhereUserIdScope(userId uint) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		fmt.Println("user id is ", userId)
		return db.Where("user_id = ?", userId)
	}
}

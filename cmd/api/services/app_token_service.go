package services

import (
	"budget-app/internal/model"
	"errors"
	"gorm.io/gorm"
	"math/rand"
	"strconv"
	"time"
)

type AppTokenService struct {
	db *gorm.DB
}

func NewAppTokenService(db *gorm.DB) *AppTokenService {
	return &AppTokenService{db: db}
}

func (appTokenService *AppTokenService) getToken() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(99999-10000+1) + 10000
}

func (appTokenService *AppTokenService) GenerateResetPasswordToken(user model.User) (*model.AppToken, error) {
	token := model.AppToken{
		TargetId:  user.ID,
		Type:      "reset_password",
		Token:     strconv.Itoa(appTokenService.getToken()),
		Used:      false,
		ExpiresAt: time.Now().Add(time.Hour * 24),
	}
	result := appTokenService.db.Create(&token)
	if result.Error != nil {
		return nil, result.Error
	}
	return &token, nil
}

func (appTokenService *AppTokenService) ValidateResetPasswordToken(user model.User, token string) (*model.AppToken, error) {
	var retrievedToken model.AppToken
	result := appTokenService.db.Where(&model.AppToken{
		TargetId: user.ID,
		Token:    token}).First(&retrievedToken)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("Invalid password reset token")
		}
		return nil, result.Error
	}

	if retrievedToken.Used {
		return nil, errors.New("Invalid password reset token")
	}

	if retrievedToken.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("Token is expired")
	}

	return &retrievedToken, nil
}

func (appTokenService *AppTokenService) InvalidateToken(userId uint, appToken model.AppToken) {
	appTokenService.db.Model(&model.AppToken{}).Where("target_id = ? AND token = ?", userId, appToken.Token).Update("used", true)
}

package services

import (
	"budget-app/cmd/api/request"
	"budget-app/common"
	"budget-app/internal/model"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (userService *UserService) CreateUser(userRequest request.RegisterUserRequest) (*model.User, error) {
	password, err := common.HashPassword(userRequest.Password)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("User registration failed")
	}
	userModel := model.User{
		FirstName: &userRequest.Firstname,
		LastName:  &userRequest.Lastname,
		Email:     userRequest.Email,
		Password:  password,
	}
	result := userService.db.Create(&userModel)
	if result.Error != nil {
		fmt.Println(result.Error)
		return nil, errors.New("User registration failed")
	}
	return &userModel, nil
}

func (userService *UserService) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	result := userService.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (userService *UserService) ChangePassword(newPassword string, user model.User) error {
	hashedPassword, err := common.HashPassword(newPassword)
	if err != nil {
		fmt.Println(err)
		return errors.New("Password changing failed")
	}
	result := userService.db.Model(user).Update("password", hashedPassword)
	if result.Error != nil {
		fmt.Println(result.Error)
		return errors.New("Password changing failed")
	}
	return nil

}

package repository

import (
	"greenenvironment/features/users"
	"greenenvironment/constant"
	"greenenvironment/helper"

	"gorm.io/gorm"
	"github.com/google/uuid"
)

type UserData struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) users.UserRepoInterface {
	return &UserData{
		DB: db,
	}
}

func (u *UserData) Register(newUser users.User) (users.User, error) {
	isEmailExist := u.IsEmailExist(newUser.Email)
	if isEmailExist {
		return users.User{}, constant.ErrEmailAlreadyExist
	}
	isUsernameExist := u.IsUsernameExist(newUser.Username)
	if isUsernameExist {
		return users.User{}, constant.ErrUsernameAlreadyExist
	}

	newUsers := User{
		ID:        uuid.New().String(),
		Name:      newUser.Name,
		Email:     newUser.Email,
		Username:  newUser.Username,
		Password:  newUser.Password,
		Coin:      0,
		Exp:       0,
		AvatarURL: "",
	}

	if err := u.DB.Create(&newUsers).Error; err != nil {
		return users.User{}, constant.ErrRegister
	}

	return users.User{
		ID:       newUsers.ID,
		Name:     newUsers.Name,
		Email:    newUsers.Email,
		Username: newUsers.Username,
		Coin:     newUsers.Coin,
		Exp:      newUsers.Exp,
	}, nil
}


func (u *UserData) Login(user users.User) (users.User, error) {
	var UserLoginData User
	result := u.DB.Where("email = ?", user.Email).First(&UserLoginData)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return users.User{}, constant.UserNotFound
		}
		return users.User{}, result.Error
	}

	if !helper.CheckPasswordHash(user.Password, UserLoginData.Password) {
		return users.User{}, constant.ErrLoginIncorrectPassword
	}
	var userLogin users.User
	userLogin.ID = UserLoginData.ID
	userLogin.Name = UserLoginData.Name
	userLogin.Email = UserLoginData.Email
	userLogin.Username = UserLoginData.Username
	userLogin.Address = UserLoginData.Address
	return userLogin, nil
}

func (u *UserData) IsUsernameExist(username string) bool {
	var user User
	if err := u.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return false
	}
	return true
}

func (u *UserData) IsEmailExist(email string) bool {
	var user User
	if err := u.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return false
	}
	return true
}
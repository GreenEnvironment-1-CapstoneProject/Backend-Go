package repository

import (
	"greenenvironment/constant"
	"greenenvironment/features/users"
	"greenenvironment/helper"

	"github.com/google/uuid"
	"gorm.io/gorm"
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

func (u *UserData) Update(user users.UserUpdate) (users.User, error) {
	var existingUser users.User
	err := u.DB.Where("id = ?", user.ID).First(&existingUser).Error
	if err != nil {
		return users.User{}, err
	}

	if user.Email != "" && user.Email != existingUser.Email {
		var count int64
		u.DB.Table("users").Where("email = ?", user.Email).Count(&count)
		if count > 0 {
			return users.User{}, constant.ErrEmailAlreadyExist
		}
	}

	if user.Username != "" && user.Username != existingUser.Username {
		var count int64
		u.DB.Table("users").Where("username = ?", user.Username).Count(&count)
		if count > 0 {
			return users.User{}, constant.ErrUsernameAlreadyExist
		}
	}

	if err := u.DB.Model(&User{}).Where("id = ?", user.ID).Updates(user).Error; err != nil {
		return users.User{}, constant.ErrUpdateUser
	}

	var updatedUser users.User
	updatedUser, err = u.GetUserByID(user.ID)
	if err != nil {
		return users.User{}, err
	}
	return updatedUser, nil
}

func (u *UserData) Delete(user users.User) error {
	_, err := u.GetUserByID(user.ID)
	if err != nil {
		return err
	}
	if err := u.DB.Where("id = ?", user.ID).Delete(&user).Error; err != nil {
		return constant.ErrDeleteUser
	}
	return nil
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

func (u *UserData) GetUserByID(id string) (users.User, error) {
	var user users.User
	var count int64
	u.DB.Table("users").Where("id = ?", id).Count(&count)
	if count == 0 {
		return users.User{}, constant.UserNotFound
	}
	if err := u.DB.Where("id = ?", id).First(&user).Error; err != nil {
		return users.User{}, constant.UserNotFound
	}
	return user, nil
}

// Admin
func (u *UserData) GetUserByIDForAdmin(id string) (users.User, error) {
	var users users.User
	res := u.DB.Model(&User{}).Where("id = ? AND deleted_at IS NULL", id).First(&users)
	if res.Error != nil {
		return users, res.Error
	}
	return users, nil
}

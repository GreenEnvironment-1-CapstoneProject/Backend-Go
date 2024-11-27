package service

import (
	"greenenvironment/constant"
	"greenenvironment/features/users"
	"greenenvironment/helper"
	"strings"
)

type UserService struct {
	userRepo users.UserRepoInterface
	jwt      helper.JWTInterface
}

func NewUserService(data users.UserRepoInterface, jwt helper.JWTInterface) users.UserServiceInterface {
	return &UserService{
		userRepo: data,
		jwt:      jwt,
	}
}

func (s *UserService) Register(user users.User) (users.User, error) {
	switch {
	case user.Email == "":
		return users.User{}, constant.ErrEmptyEmailRegister
	case user.Password == "":
		return users.User{}, constant.ErrEmptyPasswordRegister
	case user.Name == "":
		return users.User{}, constant.ErrEmptyNameRegister
	}

	user.Email = strings.ToLower(user.Email)

	isEmailValid := helper.ValidateEmail(user.Email)
	if !isEmailValid {
		return users.User{}, constant.ErrInvalidEmail
	}

	hashedPassword, err := helper.HashPassword(user.Password)
	if err != nil {
		return users.User{}, err
	}
	user.Password = hashedPassword

	user.Username = "user_" + helper.GenerateRandomString(8)

	createdUser, err := s.userRepo.Register(user)
	if err != nil {
		return users.User{}, err
	}

	return createdUser, nil
}


func (s *UserService) Login(user users.User) (users.UserLogin, error) {
	if user.Email == "" || user.Password == "" {
		return users.UserLogin{}, constant.ErrEmptyLogin
	}
	isEmailValid := helper.ValidateEmail(user.Email)
	if !isEmailValid {
		return users.UserLogin{}, constant.ErrInvalidEmail
	}
	user.Email = strings.ToLower(user.Email)

	userData, err := s.userRepo.Login(user)
	if err != nil {
		return users.UserLogin{}, err
	}

	var UserLogin helper.UserJWT

	UserLogin.ID = userData.ID
	UserLogin.Name = userData.Name
	UserLogin.Email = userData.Email
	UserLogin.Username = userData.Username
	UserLogin.Address = userData.Address
	UserLogin.Role = constant.RoleUser

	token, err := s.jwt.GenerateUserJWT(UserLogin)
	if err != nil {
		return users.UserLogin{}, err
	}

	var UserLoginData users.UserLogin
	UserLoginData.Token = token

	return UserLoginData, nil
}

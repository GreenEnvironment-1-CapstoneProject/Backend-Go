package service

import (
	"greenenvironment/constant"
	"greenenvironment/features/users"
	"greenenvironment/helper"
	"strings"

	"gorm.io/gorm"
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

func (s *UserService) Update(user users.UserUpdate) (users.UserUpdate, error) {
	if user.Email == "" && user.Username == "" && user.Password == "" && user.Address == "" && user.Name == "" && user.Gender == "" && user.Phone == "" && user.AvatarURL == "" {
		return users.UserUpdate{}, constant.ErrEmptyUpdate
	}

	if user.ID == "" {
		return users.UserUpdate{}, constant.ErrUpdateUser
	}

	if user.Email != "" {
		isEmailValid := helper.ValidateEmail(user.Email)
		if !isEmailValid {
			return users.UserUpdate{}, constant.ErrInvalidEmail
		}
		user.Email = strings.ToLower(user.Email)
	}

	if user.Username != "" {
		trimmedUsername := strings.TrimSpace(user.Username)
		isUsernameValid := helper.ValidateUsername(trimmedUsername)
		if !isUsernameValid {
			return users.UserUpdate{}, constant.ErrInvalidUsername
		}
		user.Username = trimmedUsername
	}

	if user.Phone != "" {
		trimmedPhone := strings.TrimSpace(user.Phone)
		isPhoneValid := helper.ValidatePhone(trimmedPhone)
		if !isPhoneValid {
			return users.UserUpdate{}, constant.ErrInvalidPhone
		}
		user.Phone = trimmedPhone
	}

	if user.Password != "" {
		hashedPassword, err := helper.HashPassword(user.Password)
		if err != nil {
			return users.UserUpdate{}, err
		}
		user.Password = hashedPassword
	}

	if !helper.IsValidInput(user.Name) || !helper.IsValidInput(user.Address) || !helper.IsValidInput(user.Gender) {
		return users.UserUpdate{}, constant.ErrFieldData
	}

	userData, err := s.userRepo.Update(user)
	if err != nil {
		return users.UserUpdate{}, err
	}

	var UserToken helper.UserJWT
	UserToken.ID = userData.ID
	UserToken.Name = userData.Name
	UserToken.Email = userData.Email
	UserToken.Username = userData.Username
	UserToken.Address = userData.Address
	UserToken.Role = constant.RoleUser

	token, err := s.jwt.GenerateUserJWT(UserToken)
	if err != nil {
		return users.UserUpdate{}, err
	}

	user.Token = token

	return user, nil
}

func (s *UserService) GetUserData(user users.User) (users.User, error) {
	return s.userRepo.GetUserByID(user.ID)
}

func (s *UserService) Delete(user users.User) error {
	if user.ID == "" {
		return constant.ErrDeleteUser
	}
	return s.userRepo.Delete(user)
}

func (s *UserService) RegisterOrLoginGoogle(user users.User) (users.User, error) {
	existingUser, err := s.userRepo.GetUserByEmail(user.Email)
	if err != nil && err != gorm.ErrRecordNotFound {
		return users.User{}, err
	}

	if existingUser.ID != "" {
		// User exists, return existing user
		return existingUser, nil
	}

	// Register new user
	user.Username = "google_" + helper.GenerateRandomString(8)
	newUser, err := s.Register(user)
	if err != nil {
		return users.User{}, err
	}

	return newUser, nil
}

// Admin
func (s *UserService) GetUserByIDForAdmin(id string) (users.User, error) {
	if id == "" {
		return users.User{}, constant.ErrUserIDNotFound
	}
	return s.userRepo.GetUserByIDForAdmin(id)
}
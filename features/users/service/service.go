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

func (s *UserService) Update(user users.UserUpdate) error {
	if user.ID == "" {
		return constant.ErrUpdateUser
	}

	if user.Phone != "" {
		trimmedPhone := strings.TrimSpace(user.Phone)
		isPhoneValid := helper.ValidatePhone(trimmedPhone)
		if !isPhoneValid {
			return constant.ErrInvalidPhone
		}
		user.Phone = trimmedPhone
	}

	if user.Password != "" {
		hashedPassword, err := helper.HashPassword(user.Password)
		if err != nil {
			return err
		}
		user.Password = hashedPassword
	}

	_, err := s.userRepo.Update(user)
	if err != nil {
		return err
	}

	return nil
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
	newUser, err := s.userRepo.Register(user)
	if err != nil {
		return users.User{}, err
	}

	return newUser, nil
}

func (s *UserService) UpdateAvatar(userID, avatarURL string) error {
	err := s.userRepo.UpdateAvatar(userID, avatarURL)
	if err != nil {
		return err
	}
	return nil
}

// Admin
func (s *UserService) GetUserByIDForAdmin(id string) (users.User, error) {
	if id == "" {
		return users.User{}, constant.ErrUserIDNotFound
	}
	return s.userRepo.GetUserByIDForAdmin(id)
}

func (s *UserService) GetAllByPageForAdmin(page int, limit int) ([]users.User, int, error) {
	return s.userRepo.GetAllByPageForAdmin(page, limit)
}

func (s *UserService) UpdateUserForAdmin(user users.UpdateUserByAdmin) error {
	if user.ID == "" {
		return constant.ErrUserIDNotFound
	}
	return s.userRepo.UpdateUserForAdmin(user)
}

func (s *UserService) DeleteUserForAdmin(userID string) error {
	return s.userRepo.DeleteUserForAdmin(userID)
}

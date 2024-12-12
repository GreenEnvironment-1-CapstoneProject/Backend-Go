package service

import (
	"greenenvironment/configs"
	"greenenvironment/constant"
	"greenenvironment/features/users"
	"greenenvironment/helper"
	"strings"

	"github.com/google/uuid"
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

func (s *UserService) RequestRegisterOTP(name, email, password string) error {
	if email == "" || name == "" || password == "" {
		return constant.ErrInvalidInput
	}

	hashedPassword, err := helper.HashPassword(password)
	if err != nil {
		return err
	}

	tempUser := users.TemporaryUser{
		ID:       uuid.New().String(),
		Name:     name,
		Email:    email,
		Password: hashedPassword,
	}

	err = s.userRepo.SaveTemporaryUser(tempUser)
	if err != nil {
		return err
	}

	otp := helper.NewOTP().GenerateOTP()
	expiration := helper.NewOTP().OTPExpiration(5)
	err = s.userRepo.SaveOTP(email, otp, expiration)
	if err != nil {
		return err
	}

	smtpConfig := configs.InitConfig().SMTP
	subject := "Your OTP for Registration"
	body := "Your OTP is: " + otp + ". It expires in 5 minutes."
	return helper.SendEmail(smtpConfig, email, subject, body)
}

func (s *UserService) VerifyRegisterOTP(otp string) (users.User, error) {
	if otp == "" {
		return users.User{}, constant.ErrInvalidInput
	}

	verifyData, err := s.userRepo.GetVerifyOTP(otp)
	if err != nil {
		return users.User{}, constant.ErrOTPNotValid
	}

	tempUser, err := s.userRepo.GetTemporaryUserByEmail(verifyData.Email)
	if err != nil {
		return users.User{}, err
	}

	username := "user_" + helper.GenerateRandomString(8)

	newUser := users.User{
		ID:       tempUser.ID,
		Username: username,
		Name:     tempUser.Name,
		Email:    tempUser.Email,
		Password: tempUser.Password,
	}

	createdUser, err := s.userRepo.Register(newUser)
	if err != nil {
		return users.User{}, err
	}

	err = s.userRepo.DeleteTemporaryUserByEmail(tempUser.Email)
	if err != nil {
		return users.User{}, err
	}

	err = s.userRepo.DeleteVerifyOTP(otp)
	if err != nil {
		return users.User{}, err
	}

	return createdUser, nil
}

func (s *UserService) IsEmailExist(email string) bool {
	return s.userRepo.IsEmailExist(email)
}

func (s *UserService) RequestPasswordResetOTP(email string) error {
	otp := helper.NewOTP().GenerateOTP()
	expiration := helper.NewOTP().OTPExpiration(5)

	err := s.userRepo.SaveOTP(email, otp, expiration)
	if err != nil {
		return err
	}

	smtpConfig := configs.InitConfig().SMTP
	subject := "Your OTP for Password Reset"
	body := "Your OTP is: " + otp + ". It expires in 5 minutes."

	return helper.SendEmail(smtpConfig, email, subject, body)
}

func (s *UserService) VerifyPasswordResetOTP(otp string) error {
	if otp == "" {
		return constant.ErrInvalidInput
	}

	isValidOTP := s.userRepo.ValidateOTPByOTP(otp)
	if !isValidOTP {
		return constant.ErrOTPNotValid
	}

	return nil
}

func (s *UserService) ResetPassword(newPassword string) error {
	email, err := s.userRepo.GetEmailByLatestOTP()
	if err != nil {
		return err
	}

	hashedPassword, err := helper.HashPassword(newPassword)
	if err != nil {
		return err
	}

	err = s.userRepo.UpdatePassword(email, hashedPassword)
	if err != nil {
		return err
	}

	err = s.userRepo.DeleteVerifyOTPByEmail(email)
	if err != nil {
		return err
	}

	return nil
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


func (s *UserService) UpdateUserInfo(user users.UserUpdate) error {
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

	_, err := s.userRepo.UpdateUserInfo(user)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) RequestPasswordUpdateOTP(email string) error {
	if email == "" {
		return constant.ErrEmptyEmail
	}

	otp := helper.NewOTP().GenerateOTP()
	expiration := helper.NewOTP().OTPExpiration(5)

	err := s.userRepo.SaveOTP(email, otp, expiration)
	if err != nil {
		return err
	}

	smtpConfig := configs.InitConfig().SMTP
	subject := "Your OTP for Password Update"
	body := "Your OTP is: " + otp + ". It expires in 5 minutes."

	return helper.SendEmail(smtpConfig, email, subject, body)
}

func (s *UserService) UpdatePassword(update users.PasswordUpdate) error {
	if update.Email == "" || update.OTP == "" || update.OldPassword == "" || update.NewPassword == "" {
		return constant.ErrInvalidInput
	}

	isValidOTP := s.userRepo.ValidateOTP(update.Email, update.OTP)
	if !isValidOTP {
		return constant.ErrOTPNotValid
	}

	existingUser, err := s.userRepo.GetUserByEmail(update.Email)
	if err != nil {
		return err
	}

	isOldPasswordValid := helper.CheckPasswordHash(update.OldPassword, existingUser.Password)
	if !isOldPasswordValid {
		return constant.ErrOldPasswordMismatch
	}

	hashedPassword, err := helper.HashPassword(update.NewPassword)
	if err != nil {
		return err
	}

	err = s.userRepo.UpdatePassword(update.Email, hashedPassword)
	if err != nil {
		return err
	}

	err = s.userRepo.DeleteVerifyOTP(update.OTP)
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

package service

import (
	"errors"
	"greenenvironment/constant"
	"greenenvironment/features/users"
	"greenenvironment/helper"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserData struct {
	mock.Mock
}

func (m *MockUserData) Register(newUser users.User) (users.User, error) {
	args := m.Called(newUser)
	return args.Get(0).(users.User), args.Error(1)
}

func (m *MockUserData) SaveTemporaryUser(user users.TemporaryUser) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserData) GetTemporaryUserByEmail(email string) (users.TemporaryUser, error) {
	args := m.Called(email)
	return args.Get(0).(users.TemporaryUser), args.Error(1)
}

func (m *MockUserData) DeleteTemporaryUserByEmail(email string) error {
	args := m.Called(email)
	return args.Error(0)
}

func (m *MockUserData) GetVerifyOTP(otp string) (users.VerifyOTP, error) {
	args := m.Called(otp)
	return args.Get(0).(users.VerifyOTP), args.Error(1)
}

func (m *MockUserData) DeleteVerifyOTP(otp string) error {
	args := m.Called(otp)
	return args.Error(0)
}

func (m *MockUserData) ValidateOTPByOTP(otp string) bool {
	args := m.Called(otp)
	return args.Bool(0)
}

func (m *MockUserData) GetEmailByLatestOTP() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *MockUserData) DeleteVerifyOTPByEmail(email string) error {
	args := m.Called(email)
	return args.Error(0)
}

func (m *MockUserData) Login(user users.User) (users.User, error) {
	args := m.Called(user)
	return args.Get(0).(users.User), args.Error(1)
}

func (m *MockUserData) UpdateUserInfo(user users.UserUpdate) (users.User, error) {
	args := m.Called(user)
	return args.Get(0).(users.User), args.Error(1)
}

func (m *MockUserData) Delete(user users.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserData) IsUsernameExist(username string) bool {
	args := m.Called(username)
	return args.Bool(0)
}

func (m *MockUserData) IsEmailExist(email string) bool {
	args := m.Called(email)
	return args.Bool(0)
}

func (m *MockUserData) GetUserByID(id string) (users.User, error) {
	args := m.Called(id)
	return args.Get(0).(users.User), args.Error(1)
}

func (m *MockUserData) GetUserByEmail(email string) (users.User, error) {
	args := m.Called(email)
	return args.Get(0).(users.User), args.Error(1)
}

func (m *MockUserData) UpdateAvatar(userID, avatarURL string) error {
	args := m.Called(userID, avatarURL)
	return args.Error(0)
}

func (m *MockUserData) SaveOTP(email, otp string, expiration time.Time) error {
	args := m.Called(email, otp, expiration)
	return args.Error(0)
}

func (m *MockUserData) ValidateOTP(email, otp string) bool {
	args := m.Called(email, otp)
	return args.Bool(0)
}

func (m *MockUserData) UpdatePassword(email, hashedPassword string) error {
	args := m.Called(email, hashedPassword)
	return args.Error(0)
}

func (m *MockUserData) GetUserByIDForAdmin(id string) (users.User, error) {
	args := m.Called(id)
	return args.Get(0).(users.User), args.Error(1)
}

func (m *MockUserData) GetAllUsersForAdmin() ([]users.User, error) {
	args := m.Called()
	return args.Get(0).([]users.User), args.Error(1)
}

func (m *MockUserData) GetAllByPageForAdmin(page int, limit int) ([]users.User, int, error) {
	args := m.Called(page, limit)
	return args.Get(0).([]users.User), args.Int(1), args.Error(2)
}

func (m *MockUserData) UpdateUserForAdmin(user users.UpdateUserByAdmin) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserData) DeleteUserForAdmin(userID string) error {
	args := m.Called(userID)
	return args.Error(0)
}

func TestRequestRegisterOTP(t *testing.T) {
	mockUserRepo := new(MockUserData)
	mockJWT := new(helper.MockJWTInterface)
	mockMailer := new(helper.MockMailerInterface)
	mockOTP := new(helper.MockOTPInterface)
	mockHelper := new(helper.MockHelperInterface)

	service := NewUserService(mockUserRepo, mockJWT, mockMailer, mockOTP)

	t.Run("Success Case", func(t *testing.T) {
		name := "John Doe"
		email := "john@example.com"
		password := "password123"
		otp := "123456"
		hashedPassword := "hashed-password"
		expiration := time.Now().Add(5 * time.Minute)

		// Mock password hashing
		mockHelper.On("HashPassword", password).Return(hashedPassword, nil)

		// Mock OTP generation
		mockOTP.On("GenerateOTP").Return(otp)
		mockOTP.On("OTPExpiration", 5).Return(expiration)

		// Mock repository calls
		mockUserRepo.On("SaveTemporaryUser", mock.MatchedBy(func(u users.TemporaryUser) bool {
			return u.Name == name && u.Email == email && u.Password == hashedPassword
		})).Return(nil)
		mockUserRepo.On("SaveOTP", email, otp, expiration).Return(nil)

		// Mock mailer
		mockMailer.On("Send", email, otp, "Register Account").Return(nil)

		err := service.RequestRegisterOTP(name, email, password)

		assert.NoError(t, err)
		mockUserRepo.AssertExpectations(t)
		mockMailer.AssertExpectations(t)
		mockOTP.AssertExpectations(t)
		mockHelper.AssertExpectations(t)
	})

	t.Run("Empty Input", func(t *testing.T) {
		err := service.RequestRegisterOTP("", "", "")
		assert.Error(t, err)
		assert.Equal(t, constant.ErrInvalidInput, err)
	})

	t.Run("Password Hashing Error", func(t *testing.T) {
		name := "John Doe"
		email := "john@example.com"
		password := "password123"
		expectedErr := errors.New("hashing error")

		mockHelper.On("HashPassword", password).Return("", expectedErr)

		err := service.RequestRegisterOTP(name, email, password)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		mockHelper.AssertExpectations(t)
	})

	t.Run("Save Temporary User Error", func(t *testing.T) {
		name := "John Doe"
		email := "john@example.com"
		password := "password123"
		hashedPassword := "hashed-password"
		expectedErr := errors.New("database error")

		mockHelper.On("HashPassword", password).Return(hashedPassword, nil)
		mockUserRepo.On("SaveTemporaryUser", mock.MatchedBy(func(u users.TemporaryUser) bool {
			return u.Name == name && u.Email == email && u.Password == hashedPassword
		})).Return(expectedErr)

		err := service.RequestRegisterOTP(name, email, password)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		mockHelper.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Save OTP Error", func(t *testing.T) {
		name := "John Doe"
		email := "john@example.com"
		password := "password123"
		hashedPassword := "hashed-password"
		otp := "123456"
		expiration := time.Now().Add(5 * time.Minute)
		expectedErr := errors.New("otp error")

		mockHelper.On("HashPassword", password).Return(hashedPassword, nil)
		mockUserRepo.On("SaveTemporaryUser", mock.MatchedBy(func(u users.TemporaryUser) bool {
			return u.Name == name && u.Email == email && u.Password == hashedPassword
		})).Return(nil)

		mockOTP.On("GenerateOTP").Return(otp)
		mockOTP.On("OTPExpiration", 5).Return(expiration)
		mockUserRepo.On("SaveOTP", email, otp, expiration).Return(expectedErr)

		err := service.RequestRegisterOTP(name, email, password)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		mockHelper.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
		mockOTP.AssertExpectations(t)
	})

	t.Run("Send Email Error", func(t *testing.T) {
		name := "John Doe"
		email := "john@example.com"
		password := "password123"
		hashedPassword := "hashed-password"
		otp := "123456"
		expiration := time.Now().Add(5 * time.Minute)
		expectedErr := errors.New("email error")

		mockHelper.On("HashPassword", password).Return(hashedPassword, nil)
		mockUserRepo.On("SaveTemporaryUser", mock.MatchedBy(func(u users.TemporaryUser) bool {
			return u.Name == name && u.Email == email && u.Password == hashedPassword
		})).Return(nil)

		mockOTP.On("GenerateOTP").Return(otp)
		mockOTP.On("OTPExpiration", 5).Return(expiration)
		mockUserRepo.On("SaveOTP", email, otp, expiration).Return(nil)
		mockMailer.On("Send", email, otp, "Register Account").Return(expectedErr)

		err := service.RequestRegisterOTP(name, email, password)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		mockHelper.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
		mockMailer.AssertExpectations(t)
		mockOTP.AssertExpectations(t)
	})
}

// func TestRequestRegisterOTP(t *testing.T) {
// 	mockUserRepo := new(MockUserData)
// 	mockMailer := new(helper.MockMailer)

// 	service := NewUserService(mockUserRepo, nil)

// 	t.Run("Success Case", func(t *testing.T) {
// 		name := "John Doe"
// 		email := "john.doe@example.com"
// 		password := "password123"

// 		otp := "123456"
// 		expiration := helper.NewOTP().OTPExpiration(5)

// 		mockUserRepo.On("SaveTemporaryUser", mock.MatchedBy(func(u users.TemporaryUser) bool {
// 			match := helper.CheckPasswordHash(password, u.Password)
// 			return u.Name == name && u.Email == email && match
// 		})).Return(nil).Once()

// 		mockUserRepo.On("SaveOTP", email, otp, expiration).Return(nil).Once()
// 		mockMailer.On("Send", email, otp, "Register Account").Return(nil).Once()

// 		err := service.RequestRegisterOTP(name, email, password)
// 		assert.Nil(t, err)

// 		mockUserRepo.AssertExpectations(t)
// 		mockMailer.AssertExpectations(t)
// 	})

// 	t.Run("Missing Input", func(t *testing.T) {
// 		name := ""
// 		email := ""
// 		password := ""

// 		err := service.RequestRegisterOTP(name, email, password)
// 		assert.Equal(t, constant.ErrInvalidInput, err)
// 	})

// 	t.Run("Error Saving Temporary User", func(t *testing.T) {
// 		name := "John Doe"
// 		email := "john.doe@example.com"
// 		password := "password123"
// 		hashedPassword, _ := helper.HashPassword(password)

// 		tempUser := users.TemporaryUser{
// 			ID:       "some-uuid",
// 			Name:     name,
// 			Email:    email,
// 			Password: hashedPassword,
// 		}

// 		mockUserRepo.On("SaveTemporaryUser", tempUser).Return(errors.New("failed to save user")).Once()

// 		err := service.RequestRegisterOTP(name, email, password)
// 		assert.NotNil(t, err)
// 		assert.Equal(t, "failed to save user", err.Error())
// 		mockUserRepo.AssertExpectations(t)
// 	})

// 	t.Run("Error Sending OTP", func(t *testing.T) {
// 		name := "John Doe"
// 		email := "john.doe@example.com"
// 		password := "password123"
// 		hashedPassword, _ := helper.HashPassword(password)

// 		tempUser := users.TemporaryUser{
// 			ID:       "some-uuid",
// 			Name:     name,
// 			Email:    email,
// 			Password: hashedPassword,
// 		}

// 		otp := "123456"
// 		expiration := helper.NewOTP().OTPExpiration(5)

// 		mockUserRepo.On("SaveTemporaryUser", tempUser).Return(nil).Once()
// 		mockUserRepo.On("SaveOTP", email, otp, expiration).Return(nil).Once()
// 		mockMailer.On("Send", email, otp, "Register Account").Return(errors.New("failed to send email")).Once()

// 		err := service.RequestRegisterOTP(name, email, password)
// 		assert.NotNil(t, err)
// 		assert.Equal(t, "failed to send email", err.Error())
// 		mockUserRepo.AssertExpectations(t)
// 		mockMailer.AssertExpectations(t)
// 	})
// }

// func TestVerifyRegisterOTP(t *testing.T) {
// 	mockRepo := new(MockUserData)
// 	mockJWT := new(helper.MockJWTInterface)
// 	service := NewUserService(mockRepo, mockJWT)

// 	otp := "123456"
// 	tempUser := users.TemporaryUser{
// 		ID:       "1",
// 		Name:     "Test User",
// 		Email:    "test@example.com",
// 		Password: "hashedpassword123",
// 	}

// 	mockRepo.On("GetVerifyOTP", otp).Return(users.VerifyOTP{Email: tempUser.Email}, nil).Once()
// 	mockRepo.On("GetTemporaryUserByEmail", tempUser.Email).Return(tempUser, nil).Once()
// 	mockRepo.On("Register", mock.Anything).Return(users.User{}, nil).Once()
// 	mockRepo.On("DeleteTemporaryUserByEmail", tempUser.Email).Return(nil).Once()
// 	mockRepo.On("DeleteVerifyOTP", otp).Return(nil).Once()

// 	_, err := service.VerifyRegisterOTP(otp)

// 	assert.NoError(t, err)
// 	mockRepo.AssertExpectations(t)
// }

// func TestRequestPasswordResetOTP(t *testing.T) {
// 	mockRepo := new(MockUserData)
// 	mockJWT := new(helper.MockJWTInterface)
// 	service := NewUserService(mockRepo, mockJWT)

// 	email := "test@example.com"
// 	otp := "123456"

// 	mockRepo.On("SaveOTP", email, otp, mock.Anything).Return(nil).Once()

// 	err := service.RequestPasswordResetOTP(email)

// 	assert.NoError(t, err)
// 	mockRepo.AssertExpectations(t)
// }

// func TestVerifyPasswordResetOTP(t *testing.T) {
// 	mockRepo := new(MockUserData)
// 	mockJWT := new(helper.MockJWTInterface)
// 	service := NewUserService(mockRepo, mockJWT)

// 	otp := "123456"

// 	mockRepo.On("ValidateOTPByOTP", otp).Return(true).Once()

// 	err := service.VerifyPasswordResetOTP(otp)

// 	assert.NoError(t, err)
// 	mockRepo.AssertExpectations(t)
// }

// func TestResetPassword(t *testing.T) {
// 	mockRepo := new(MockUserData)
// 	mockJWT := new(helper.MockJWTInterface)
// 	service := NewUserService(mockRepo, mockJWT)

// 	email := "test@example.com"
// 	newPassword := "newpassword123"
// 	hashedPassword := "hashednewpassword123"

// 	mockRepo.On("GetEmailByLatestOTP", mock.Anything).Return(email, nil).Once()
// 	mockRepo.On("UpdatePassword", email, hashedPassword).Return(nil).Once()
// 	mockRepo.On("DeleteVerifyOTPByEmail", email).Return(nil).Once()

// 	err := service.ResetPassword(newPassword)

// 	assert.NoError(t, err)
// 	mockRepo.AssertExpectations(t)
// }

// func TestLogin(t *testing.T) {
// 	mockRepo := new(MockUserData)
// 	mockJWT := new(helper.MockJWTInterface)
// 	service := NewUserService(mockRepo, mockJWT)

// 	user := users.User{
// 		Email:    "test@example.com",
// 		Password: "password123",
// 	}
// 	userData := users.User{
// 		ID:       "1",
// 		Name:     "Test User",
// 		Email:    "test@example.com",
// 		Username: "testuser",
// 		Address:  "Test Address",
// 	}
// 	mockRepo.On("Login", user).Return(userData, nil).Once()

// 	token := "testtoken"
// 	mockJWT.On("GenerateUserJWT", mock.Anything).Return(token, nil).Once()

// 	result, err := service.Login(user)

// 	assert.NoError(t, err)
// 	assert.Equal(t, token, result.Token)
// 	mockRepo.AssertExpectations(t)
// 	mockJWT.AssertExpectations(t)
// }

// func TestUpdateUserInfo(t *testing.T) {
// 	mockRepo := new(MockUserData)
// 	mockJWT := new(helper.MockJWTInterface)
// 	service := NewUserService(mockRepo, mockJWT)

// 	userUpdate := users.UserUpdate{
// 		ID:      "1",
// 		Name:    "Updated User",
// 		Phone:   "08123456789",
// 		Address: "Updated Address",
// 	}

// 	mockRepo.On("UpdateUserInfo", userUpdate).Return(users.User{}, nil).Once()

// 	err := service.UpdateUserInfo(userUpdate)

// 	assert.NoError(t, err)
// 	mockRepo.AssertExpectations(t)
// }

// func TestRequestPasswordUpdateOTP(t *testing.T) {
// 	mockRepo := new(MockUserData)
// 	mockJWT := new(helper.MockJWTInterface)
// 	service := NewUserService(mockRepo, mockJWT)

// 	email := "test@example.com"
// 	otp := "123456"

// 	mockRepo.On("SaveOTP", email, otp, mock.Anything).Return(nil).Once()

// 	err := service.RequestPasswordUpdateOTP(email)

// 	assert.NoError(t, err)
// 	mockRepo.AssertExpectations(t)
// }

// func TestUpdatePassword(t *testing.T) {
// 	mockRepo := new(MockUserData)
// 	mockJWT := new(helper.MockJWTInterface)
// 	service := NewUserService(mockRepo, mockJWT)

// 	update := users.PasswordUpdate{
// 		Email:       "test@example.com",
// 		OTP:         "123456",
// 		OldPassword: "oldpassword123",
// 		NewPassword: "newpassword123",
// 	}

// 	existingUser := users.User{
// 		Email:    "test@example.com",
// 		Password: "hashedoldpassword123",
// 	}

// 	mockRepo.On("ValidateOTP", update.Email, update.OTP).Return(true).Once()
// 	mockRepo.On("GetUserByEmail", update.Email).Return(existingUser, nil).Once()
// 	mockRepo.On("UpdatePassword", update.Email, mock.Anything).Return(nil).Once()
// 	mockRepo.On("DeleteVerifyOTP", update.OTP).Return(nil).Once()

// 	err := service.UpdatePassword(update)

// 	assert.NoError(t, err)
// 	mockRepo.AssertExpectations(t)
// }

// func TestGetUserData(t *testing.T) {
// 	mockRepo := new(MockUserData)
// 	mockJWT := new(helper.MockJWTInterface)
// 	service := NewUserService(mockRepo, mockJWT)

// 	user := users.User{
// 		ID: "1",
// 	}

// 	expectedUser := users.User{
// 		ID:   "1",
// 		Name: "Test User",
// 	}

// 	mockRepo.On("GetUserByID", user.ID).Return(expectedUser, nil).Once()

// 	result, err := service.GetUserData(user)

// 	assert.NoError(t, err)
// 	assert.Equal(t, expectedUser, result)
// 	mockRepo.AssertExpectations(t)
// }

// func TestDelete(t *testing.T) {
// 	mockRepo := new(MockUserData)
// 	mockJWT := new(helper.MockJWTInterface)
// 	service := NewUserService(mockRepo, mockJWT)

// 	user := users.User{
// 		ID: "1",
// 	}

// 	mockRepo.On("Delete", user).Return(nil).Once()

// 	err := service.Delete(user)

// 	assert.NoError(t, err)
// 	mockRepo.AssertExpectations(t)
// }

// func TestRegisterOrLoginGoogle(t *testing.T) {
// 	mockRepo := new(MockUserData)
// 	mockJWT := new(helper.MockJWTInterface)
// 	service := NewUserService(mockRepo, mockJWT)

// 	user := users.User{
// 		Email:    "test@example.com",
// 		Name:     "Google User",
// 		Username: "",
// 	}

// 	existingUser := users.User{
// 		ID:    "1",
// 		Email: "test@example.com",
// 	}

// 	mockRepo.On("GetUserByEmail", user.Email).Return(existingUser, nil).Once()

// 	result, err := service.RegisterOrLoginGoogle(user)

// 	assert.NoError(t, err)
// 	assert.Equal(t, existingUser, result)
// 	mockRepo.AssertExpectations(t)
// }

// func TestUpdateAvatar(t *testing.T) {
// 	mockRepo := new(MockUserData)
// 	mockJWT := new(helper.MockJWTInterface)
// 	service := NewUserService(mockRepo, mockJWT)

// 	userID := "1"
// 	avatarURL := "https://example.com/avatar.png"

// 	mockRepo.On("UpdateAvatar", userID, avatarURL).Return(nil).Once()

// 	err := service.UpdateAvatar(userID, avatarURL)

// 	assert.NoError(t, err)
// 	mockRepo.AssertExpectations(t)
// }

// func TestGetUserByIDForAdmin(t *testing.T) {
// 	mockRepo := new(MockUserData)
// 	mockJWT := new(helper.MockJWTInterface)
// 	service := NewUserService(mockRepo, mockJWT)

// 	userID := "1"

// 	expectedUser := users.User{
// 		ID:   "1",
// 		Name: "Admin User",
// 	}

// 	mockRepo.On("GetUserByIDForAdmin", userID).Return(expectedUser, nil).Once()

// 	result, err := service.GetUserByIDForAdmin(userID)

// 	assert.NoError(t, err)
// 	assert.Equal(t, expectedUser, result)
// 	mockRepo.AssertExpectations(t)
// }

// func TestGetAllByPageForAdmin(t *testing.T) {
// 	mockRepo := new(MockUserData)
// 	mockJWT := new(helper.MockJWTInterface)
// 	service := NewUserService(mockRepo, mockJWT)

// 	page := 1
// 	limit := 10

// 	users := []users.User{
// 		{ID: "1", Name: "User 1"},
// 		{ID: "2", Name: "User 2"},
// 	}
// 	totalPages := 1

// 	mockRepo.On("GetAllByPageForAdmin", page, limit).Return(users, totalPages, nil).Once()

// 	result, total, err := service.GetAllByPageForAdmin(page, limit)

// 	assert.NoError(t, err)
// 	assert.Equal(t, users, result)
// 	assert.Equal(t, totalPages, total)
// 	mockRepo.AssertExpectations(t)
// }

// func TestUpdateUserForAdmin(t *testing.T) {
// 	mockRepo := new(MockUserData)
// 	mockJWT := new(helper.MockJWTInterface)
// 	service := NewUserService(mockRepo, mockJWT)

// 	userUpdate := users.UpdateUserByAdmin{
// 		ID:    "1",
// 		Name:  "Updated Admin User",
// 		Address: "Updated Address",
// 		Gender: "Male",
// 		Phone: "08123456789",
// 	}

// 	mockRepo.On("UpdateUserForAdmin", userUpdate).Return(nil).Once()

// 	err := service.UpdateUserForAdmin(userUpdate)

// 	assert.NoError(t, err)
// 	mockRepo.AssertExpectations(t)
// }

// func TestDeleteUserForAdmin(t *testing.T) {
// 	mockRepo := new(MockUserData)
// 	mockJWT := new(helper.MockJWTInterface)
// 	service := NewUserService(mockRepo, mockJWT)

// 	userID := "1"

// 	mockRepo.On("DeleteUserForAdmin", userID).Return(nil).Once()

// 	err := service.DeleteUserForAdmin(userID)

// 	assert.NoError(t, err)
// 	mockRepo.AssertExpectations(t)
// }

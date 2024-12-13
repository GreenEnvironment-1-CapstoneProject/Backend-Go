package service

import (
	"errors"
	"greenenvironment/constant"
	"greenenvironment/features/users"
	"greenenvironment/helper"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
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

type MockHasher struct {
	mock.Mock
}

func (m *MockHasher) HashPassword(password string) (string, error) {
	args := m.Called(password)
	if args.Get(0) == nil {
		hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		return string(hashed), nil
	}
	return args.String(0), args.Error(1)
}

func (m *MockHasher) CheckPasswordHash(password, hash string) bool {
	args := m.Called(password, hash)
	return args.Bool(0)
}

type MockJWTInterface struct {
	mock.Mock
}

func (m *MockJWTInterface) GenerateUserJWT(user helper.UserJWT) (string, error) {
	args := m.Called(user)
	return args.String(0), args.Error(1)
}

func (m *MockJWTInterface) GenerateAdminJWT(user helper.AdminJWT) (string, error) {
	args := m.Called(user)
	return args.String(0), args.Error(1)
}

func (m *MockJWTInterface) GenerateUserToken(user helper.UserJWT) string {
	args := m.Called(user)
	return args.String(0)
}

func (m *MockJWTInterface) GenerateAdminToken(user helper.AdminJWT) string {
	args := m.Called(user)
	return args.String(0)
}

func (m *MockJWTInterface) ExtractUserToken(token *jwt.Token) map[string]interface{} {
	args := m.Called(token)
	return args.Get(0).(map[string]interface{})
}

func (m *MockJWTInterface) ExtractAdminToken(token *jwt.Token) map[string]interface{} {
	args := m.Called(token)
	return args.Get(0).(map[string]interface{})
}

func (m *MockJWTInterface) ValidateToken(token string) (*jwt.Token, error) {
	args := m.Called(token)
	return args.Get(0).(*jwt.Token), args.Error(1)
}

type MockMailerInterface struct {
	mock.Mock
}

func (m *MockMailerInterface) Send(email, otpCode, subject string) error {
	args := m.Called(email, otpCode, subject)
	return args.Error(0)
}

type MockOTPInterface struct {
	mock.Mock
}

func (m *MockOTPInterface) GenerateOTP() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockOTPInterface) OTPExpiration(durationMinutes int) time.Time {
	args := m.Called(durationMinutes)
	return args.Get(0).(time.Time)
}

type MockPasswordInterface struct {
	mock.Mock
}

func (m *MockPasswordInterface) HashPassword(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockPasswordInterface) CheckPasswordHash(hashedPassword, password string) error {
	args := m.Called(hashedPassword, password)
	return args.Error(0)
}

func TestRequestRegisterOTP(t *testing.T) {
	mockRepo := new(MockUserData)
	mockMailer := new(MockMailerInterface)
	mockOTP := new(MockOTPInterface)

	svc := NewUserService(mockRepo, nil, mockMailer, mockOTP)

	t.Run("success", func(t *testing.T) {
		mockRepo.On("SaveTemporaryUser", mock.Anything).Return(nil).Once()
		mockRepo.On("SaveOTP", "test@example.com", mock.Anything, mock.Anything).Return(nil).Once()
		mockMailer.On("Send", "test@example.com", mock.Anything, "Register Account").Return(nil).Once()
		mockOTP.On("GenerateOTP").Return("123456").Once()
		mockOTP.On("OTPExpiration", 5).Return(time.Now().Add(5 * time.Minute)).Once()

		err := svc.RequestRegisterOTP("Test User", "test@example.com", "password123")
		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)
		mockMailer.AssertExpectations(t)
		mockOTP.AssertExpectations(t)
	})

	t.Run("error invalid input", func(t *testing.T) {
		err := svc.RequestRegisterOTP("", "", "")
		assert.EqualError(t, err, constant.ErrInvalidInput.Error())
	})

	t.Run("error saving temporary user", func(t *testing.T) {
		mockRepo.On("SaveTemporaryUser", mock.Anything).Return(errors.New("failed to save temp user")).Once()

		err := svc.RequestRegisterOTP("Test User", "test@example.com", "password123")
		assert.EqualError(t, err, "failed to save temp user")

		mockRepo.AssertExpectations(t)
	})

	t.Run("error sending otp", func(t *testing.T) {
		mockRepo.On("SaveTemporaryUser", mock.Anything).Return(nil).Once()
		mockRepo.On("SaveOTP", "test@example.com", mock.Anything, mock.Anything).Return(nil).Once()
		mockMailer.On("Send", "test@example.com", mock.Anything, "Register Account").Return(errors.New("failed to send email")).Once()
		mockOTP.On("GenerateOTP").Return("123456").Once()
		mockOTP.On("OTPExpiration", 5).Return(time.Now().Add(5 * time.Minute)).Once()

		err := svc.RequestRegisterOTP("Test User", "test@example.com", "password123")
		assert.EqualError(t, err, "failed to send email")

		mockRepo.AssertExpectations(t)
		mockMailer.AssertExpectations(t)
		mockOTP.AssertExpectations(t)
	})
}

func TestVerifyRegisterOTP(t *testing.T) {
	mockRepo := new(MockUserData)
	mockOTP := new(MockOTPInterface)

	svc := NewUserService(mockRepo, nil, nil, mockOTP)

	t.Run("success", func(t *testing.T) {
		mockRepo.On("GetVerifyOTP", "123456").Return(users.VerifyOTP{Email: "test@example.com"}, nil).Once()
		mockRepo.On("GetTemporaryUserByEmail", "test@example.com").Return(users.TemporaryUser{
			ID:       "1",
			Name:     "Test User",
			Email:    "test@example.com",
			Password: "hashed_password",
		}, nil).Once()
		mockRepo.On("Register", mock.Anything).Return(users.User{
			ID:       "1",
			Username: "user_random1234",
			Name:     "Test User",
			Email:    "test@example.com",
			Password: "hashed_password",
		}, nil).Once()
		mockRepo.On("DeleteTemporaryUserByEmail", "test@example.com").Return(nil).Once()
		mockRepo.On("DeleteVerifyOTP", "123456").Return(nil).Once()

		result, err := svc.VerifyRegisterOTP("123456")
		assert.NoError(t, err)
		assert.Equal(t, "test@example.com", result.Email)

		mockRepo.AssertExpectations(t)
	})

	t.Run("error invalid input", func(t *testing.T) {
		result, err := svc.VerifyRegisterOTP("")
		assert.EqualError(t, err, constant.ErrInvalidInput.Error())
		assert.Equal(t, users.User{}, result)
	})

	t.Run("error getting verify OTP", func(t *testing.T) {
		mockRepo.On("GetVerifyOTP", "123456").Return(users.VerifyOTP{}, errors.New("OTP not valid")).Once()

		result, err := svc.VerifyRegisterOTP("123456")
		assert.EqualError(t, err, constant.ErrOTPNotValid.Error())
		assert.Equal(t, users.User{}, result)

		mockRepo.AssertExpectations(t)
	})

	t.Run("error getting temporary user", func(t *testing.T) {
		mockRepo.On("GetVerifyOTP", "123456").Return(users.VerifyOTP{Email: "test@example.com"}, nil).Once()
		mockRepo.On("GetTemporaryUserByEmail", "test@example.com").Return(users.TemporaryUser{}, errors.New("user not found")).Once()

		result, err := svc.VerifyRegisterOTP("123456")
		assert.EqualError(t, err, "user not found")
		assert.Equal(t, users.User{}, result)

		mockRepo.AssertExpectations(t)
	})

	t.Run("error registering user", func(t *testing.T) {
		mockRepo.On("GetVerifyOTP", "123456").Return(users.VerifyOTP{Email: "test@example.com"}, nil).Once()
		mockRepo.On("GetTemporaryUserByEmail", "test@example.com").Return(users.TemporaryUser{
			ID:       "1",
			Name:     "Test User",
			Email:    "test@example.com",
			Password: "hashed_password",
		}, nil).Once()
		mockRepo.On("Register", mock.Anything).Return(users.User{}, errors.New("registration failed")).Once()

		result, err := svc.VerifyRegisterOTP("123456")
		assert.EqualError(t, err, "registration failed")
		assert.Equal(t, users.User{}, result)

		mockRepo.AssertExpectations(t)
	})

	t.Run("error deleting temporary user", func(t *testing.T) {
		mockRepo.On("GetVerifyOTP", "123456").Return(users.VerifyOTP{Email: "test@example.com"}, nil).Once()
		mockRepo.On("GetTemporaryUserByEmail", "test@example.com").Return(users.TemporaryUser{
			ID:       "1",
			Name:     "Test User",
			Email:    "test@example.com",
			Password: "hashed_password",
		}, nil).Once()
		mockRepo.On("Register", mock.Anything).Return(users.User{
			ID:       "1",
			Username: "user_random1234",
			Name:     "Test User",
			Email:    "test@example.com",
			Password: "hashed_password",
		}, nil).Once()
		mockRepo.On("DeleteTemporaryUserByEmail", "test@example.com").Return(errors.New("failed to delete temp user")).Once()

		result, err := svc.VerifyRegisterOTP("123456")
		assert.EqualError(t, err, "failed to delete temp user")
		assert.Equal(t, users.User{}, result)

		mockRepo.AssertExpectations(t)
	})

	t.Run("error deleting verify OTP", func(t *testing.T) {
		mockRepo.On("GetVerifyOTP", "123456").Return(users.VerifyOTP{Email: "test@example.com"}, nil).Once()
		mockRepo.On("GetTemporaryUserByEmail", "test@example.com").Return(users.TemporaryUser{
			ID:       "1",
			Name:     "Test User",
			Email:    "test@example.com",
			Password: "hashed_password",
		}, nil).Once()
		mockRepo.On("Register", mock.Anything).Return(users.User{
			ID:       "1",
			Username: "user_random1234",
			Name:     "Test User",
			Email:    "test@example.com",
			Password: "hashed_password",
		}, nil).Once()
		mockRepo.On("DeleteTemporaryUserByEmail", "test@example.com").Return(nil).Once()
		mockRepo.On("DeleteVerifyOTP", "123456").Return(errors.New("failed to delete verify OTP")).Once()

		result, err := svc.VerifyRegisterOTP("123456")
		assert.EqualError(t, err, "failed to delete verify OTP")
		assert.Equal(t, users.User{}, result)

		mockRepo.AssertExpectations(t)
	})
}

func TestIsEmailExist(t *testing.T) {
	mockRepo := new(MockUserData)

	svc := NewUserService(mockRepo, nil, nil, nil)

	t.Run("email exists", func(t *testing.T) {
		mockRepo.On("IsEmailExist", "test@example.com").Return(true).Once()

		result := svc.IsEmailExist("test@example.com")
		assert.True(t, result)

		mockRepo.AssertExpectations(t)
	})

	t.Run("email does not exist", func(t *testing.T) {
		mockRepo.On("IsEmailExist", "test@example.com").Return(false).Once()

		result := svc.IsEmailExist("test@example.com")
		assert.False(t, result)

		mockRepo.AssertExpectations(t)
	})
}

func TestRequestPasswordResetOTP(t *testing.T) {
	mockRepo := new(MockUserData)
	mockMailer := new(MockMailerInterface)
	mockOTP := new(MockOTPInterface)

	svc := NewUserService(mockRepo, nil, mockMailer, mockOTP)

	t.Run("success", func(t *testing.T) {
		mockOTP.On("GenerateOTP").Return("123456").Once()
		mockOTP.On("OTPExpiration", 5).Return(time.Now().Add(5 * time.Minute)).Once()
		mockRepo.On("SaveOTP", "test@example.com", "123456", mock.Anything).Return(nil).Once()
		mockMailer.On("Send", "test@example.com", "123456", "Reset Password").Return(nil).Once()

		err := svc.RequestPasswordResetOTP("test@example.com")
		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)
		mockMailer.AssertExpectations(t)
		mockOTP.AssertExpectations(t)
	})

	t.Run("error saving OTP", func(t *testing.T) {
		mockOTP.On("GenerateOTP").Return("123456").Once()
		mockOTP.On("OTPExpiration", 5).Return(time.Now().Add(5 * time.Minute)).Once()
		mockRepo.On("SaveOTP", "test@example.com", "123456", mock.Anything).Return(errors.New("failed to save OTP")).Once()

		err := svc.RequestPasswordResetOTP("test@example.com")
		assert.EqualError(t, err, "failed to save OTP")

		mockRepo.AssertExpectations(t)
		mockOTP.AssertExpectations(t)
	})

	t.Run("error sending OTP", func(t *testing.T) {
		mockOTP.On("GenerateOTP").Return("123456").Once()
		mockOTP.On("OTPExpiration", 5).Return(time.Now().Add(5 * time.Minute)).Once()
		mockRepo.On("SaveOTP", "test@example.com", "123456", mock.Anything).Return(nil).Once()
		mockMailer.On("Send", "test@example.com", "123456", "Reset Password").Return(errors.New("failed to send email")).Once()

		err := svc.RequestPasswordResetOTP("test@example.com")
		assert.EqualError(t, err, "failed to send email")

		mockRepo.AssertExpectations(t)
		mockMailer.AssertExpectations(t)
		mockOTP.AssertExpectations(t)
	})
}

func TestVerifyPasswordResetOTP(t *testing.T) {
	mockRepo := new(MockUserData)

	svc := NewUserService(mockRepo, nil, nil, nil)

	t.Run("valid OTP", func(t *testing.T) {
		mockRepo.On("ValidateOTPByOTP", "123456").Return(true).Once()

		err := svc.VerifyPasswordResetOTP("123456")
		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid OTP", func(t *testing.T) {
		mockRepo.On("ValidateOTPByOTP", "123456").Return(false).Once()

		err := svc.VerifyPasswordResetOTP("123456")
		assert.EqualError(t, err, constant.ErrOTPNotValid.Error())

		mockRepo.AssertExpectations(t)
	})

	t.Run("error invalid input", func(t *testing.T) {
		err := svc.VerifyPasswordResetOTP("")
		assert.EqualError(t, err, constant.ErrInvalidInput.Error())
	})
}

func TestUserService_Login(t *testing.T) {
	mockUserRepo := new(MockUserData)
	mockJWT := new(MockJWTInterface)

	mockUser := users.User{
		ID:       "1",
		Name:     "John Doe",
		Email:    "johndoe@gmail.com",
		Username: "johndoe",
		Address:  "123 Street",
		Password: "password",
	}

	mockUserLogin := helper.UserJWT{
		ID:       mockUser.ID,
		Name:     mockUser.Name,
		Email:    mockUser.Email,
		Username: mockUser.Username,
		Address:  mockUser.Address,
		Role:     constant.RoleUser,
	}

	mockToken := "sample.jwt.token"

	service := &UserService{
		userRepo: mockUserRepo,
		jwt:      mockJWT,
	}

	t.Run("success", func(t *testing.T) {
		mockUserRepo.On("Login", mock.AnythingOfType("users.User")).Return(mockUser, nil).Once()
		mockJWT.On("GenerateUserJWT", mockUserLogin).Return(mockToken, nil).Once()

		result, err := service.Login(mockUser)

		assert.NoError(t, err)
		assert.Equal(t, mockToken, result.Token)

		mockUserRepo.AssertExpectations(t)
		mockJWT.AssertExpectations(t)
	})

	t.Run("error in repository", func(t *testing.T) {
		mockUserRepo.On("Login", mock.AnythingOfType("users.User")).Return(users.User{}, errors.New("repository error")).Once()

		result, err := service.Login(mockUser)

		assert.Error(t, err)
		assert.Empty(t, result.Token)

		mockUserRepo.AssertExpectations(t)
	})

	t.Run("error in JWT generation", func(t *testing.T) {
		mockUserRepo.On("Login", mock.AnythingOfType("users.User")).Return(mockUser, nil).Once()
		mockJWT.On("GenerateUserJWT", mockUserLogin).Return("", errors.New("JWT error")).Once()

		result, err := service.Login(mockUser)

		assert.Error(t, err)
		assert.Empty(t, result.Token)

		mockUserRepo.AssertExpectations(t)
		mockJWT.AssertExpectations(t)
	})
}

func TestUserService_UpdateUserInfo(t *testing.T) {
	mockUserRepo := new(MockUserData)
	service := &UserService{
		userRepo: mockUserRepo,
	}

	mockUserUpdate := users.UserUpdate{
		ID:    "1",
		Phone: " 08123456789 ",
	}

	t.Run("success", func(t *testing.T) {
		mockUserRepo.On("UpdateUserInfo", mock.AnythingOfType("users.UserUpdate")).Return(users.User{}, nil).Once()

		err := service.UpdateUserInfo(mockUserUpdate)

		assert.NoError(t, err)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("missing user ID", func(t *testing.T) {
		mockUserUpdate.ID = ""

		err := service.UpdateUserInfo(mockUserUpdate)

		assert.Error(t, err)
		assert.Equal(t, constant.ErrUpdateUser, err)
	})

	t.Run("invalid phone", func(t *testing.T) {
		mockUserUpdate.ID = "1"
		mockUserUpdate.Phone = "invalid-phone"

		err := service.UpdateUserInfo(mockUserUpdate)

		assert.Error(t, err)
		assert.Equal(t, constant.ErrInvalidPhone, err)
	})
}

func TestUserService_RequestPasswordUpdateOTP(t *testing.T) {
	mockUserRepo := new(MockUserData)
	mockOTP := new(MockOTPInterface)
	mockMailer := new(MockMailerInterface)

	service := &UserService{
		userRepo: mockUserRepo,
		otp:      mockOTP,
		mailer:   mockMailer,
	}

	mockEmail := "johndoe@gmail.com"
	mockOTPCode := "123456"
	mockExpiration := time.Now().Add(5 * time.Minute)

	t.Run("success", func(t *testing.T) {
		mockOTP.On("GenerateOTP").Return(mockOTPCode).Once()
		mockOTP.On("OTPExpiration", 5).Return(mockExpiration).Once()
		mockUserRepo.On("SaveOTP", mockEmail, mockOTPCode, mockExpiration).Return(nil).Once()
		mockMailer.On("Send", mockEmail, mockOTPCode, "Update Password").Return(nil).Once()

		err := service.RequestPasswordUpdateOTP(mockEmail)

		assert.NoError(t, err)
		mockOTP.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
		mockMailer.AssertExpectations(t)
	})

	t.Run("empty email", func(t *testing.T) {
		err := service.RequestPasswordUpdateOTP("")

		assert.Error(t, err)
		assert.Equal(t, constant.ErrEmptyEmail, err)
	})

	t.Run("error saving OTP", func(t *testing.T) {
		mockOTP.On("GenerateOTP").Return(mockOTPCode).Once()
		mockOTP.On("OTPExpiration", 5).Return(mockExpiration).Once()
		mockUserRepo.On("SaveOTP", mockEmail, mockOTPCode, mockExpiration).Return(errors.New("save OTP error")).Once()

		err := service.RequestPasswordUpdateOTP(mockEmail)

		assert.Error(t, err)
		mockOTP.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("error sending email", func(t *testing.T) {
		mockOTP.On("GenerateOTP").Return(mockOTPCode).Once()
		mockOTP.On("OTPExpiration", 5).Return(mockExpiration).Once()
		mockUserRepo.On("SaveOTP", mockEmail, mockOTPCode, mockExpiration).Return(nil).Once()
		mockMailer.On("Send", mockEmail, mockOTPCode, "Update Password").Return(errors.New("email error")).Once()

		err := service.RequestPasswordUpdateOTP(mockEmail)

		assert.Error(t, err)
		mockOTP.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
		mockMailer.AssertExpectations(t)
	})
}

func TestUserService_GetUserData(t *testing.T) {
	mockUserRepo := new(MockUserData)
	service := &UserService{
		userRepo: mockUserRepo,
	}

	mockUser := users.User{
		ID: "1",
		Name: "John Doe",
	}

	t.Run("success", func(t *testing.T) {
		mockUserRepo.On("GetUserByID", mockUser.ID).Return(mockUser, nil).Once()

		result, err := service.GetUserData(mockUser)

		assert.NoError(t, err)
		assert.Equal(t, mockUser, result)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockUserRepo.On("GetUserByID", mockUser.ID).Return(users.User{}, errors.New("repository error")).Once()

		result, err := service.GetUserData(mockUser)

		assert.Error(t, err)
		assert.Equal(t, users.User{}, result)
		mockUserRepo.AssertExpectations(t)
	})
}

func TestUserService_Delete(t *testing.T) {
	mockUserRepo := new(MockUserData)
	service := &UserService{
		userRepo: mockUserRepo,
	}

	mockUser := users.User{
		ID: "1",
	}

	t.Run("success", func(t *testing.T) {
		mockUserRepo.On("Delete", mockUser).Return(nil).Once()

		err := service.Delete(mockUser)

		assert.NoError(t, err)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("missing user ID", func(t *testing.T) {
		mockUser.ID = ""

		err := service.Delete(mockUser)

		assert.Error(t, err)
		assert.Equal(t, constant.ErrDeleteUser, err)
	})
}

func TestUserService_GetUserByIDForAdmin(t *testing.T) {
	mockUserRepo := new(MockUserData)
	service := &UserService{
		userRepo: mockUserRepo,
	}

	mockUser := users.User{
		ID: "1",
		Name: "Admin User",
	}

	t.Run("success", func(t *testing.T) {
		mockUserRepo.On("GetUserByIDForAdmin", mockUser.ID).Return(mockUser, nil).Once()

		result, err := service.GetUserByIDForAdmin(mockUser.ID)

		assert.NoError(t, err)
		assert.Equal(t, mockUser, result)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("missing user ID", func(t *testing.T) {
		result, err := service.GetUserByIDForAdmin("")

		assert.Error(t, err)
		assert.Equal(t, constant.ErrUserIDNotFound, err)
		assert.Equal(t, users.User{}, result)
	})
}

func TestUserService_GetAllByPageForAdmin(t *testing.T) {
	mockUserRepo := new(MockUserData)
	service := &UserService{
		userRepo: mockUserRepo,
	}

	mockUsers := []users.User{
		{ID: "1", Name: "User 1"},
		{ID: "2", Name: "User 2"},
	}

	mockPage := 1
	mockLimit := 10
	mockTotal := len(mockUsers)

	t.Run("success", func(t *testing.T) {
		mockUserRepo.On("GetAllByPageForAdmin", mockPage, mockLimit).Return(mockUsers, mockTotal, nil).Once()

		result, total, err := service.GetAllByPageForAdmin(mockPage, mockLimit)

		assert.NoError(t, err)
		assert.Equal(t, mockUsers, result)
		assert.Equal(t, mockTotal, total)
		mockUserRepo.AssertExpectations(t)
	})
}

func TestUserService_UpdateUserForAdmin(t *testing.T) {
	mockUserRepo := new(MockUserData)
	service := &UserService{
		userRepo: mockUserRepo,
	}

	mockUser := users.UpdateUserByAdmin{
		ID:   "1",
		Name: "Updated User",
	}

	t.Run("success", func(t *testing.T) {
		mockUserRepo.On("UpdateUserForAdmin", mockUser).Return(nil).Once()

		err := service.UpdateUserForAdmin(mockUser)

		assert.NoError(t, err)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("missing user ID", func(t *testing.T) {
		mockUser.ID = ""

		err := service.UpdateUserForAdmin(mockUser)

		assert.Error(t, err)
		assert.Equal(t, constant.ErrUserIDNotFound, err)
	})
}

func TestUserService_DeleteUserForAdmin(t *testing.T) {
	mockUserRepo := new(MockUserData)
	service := &UserService{
		userRepo: mockUserRepo,
	}

	mockUserID := "1"

	t.Run("success", func(t *testing.T) {
		mockUserRepo.On("DeleteUserForAdmin", mockUserID).Return(nil).Once()

		err := service.DeleteUserForAdmin(mockUserID)

		assert.NoError(t, err)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("error in repository", func(t *testing.T) {
		mockUserRepo.On("DeleteUserForAdmin", mockUserID).Return(errors.New("repository error")).Once()

		err := service.DeleteUserForAdmin(mockUserID)

		assert.Error(t, err)
		mockUserRepo.AssertExpectations(t)
	})
}

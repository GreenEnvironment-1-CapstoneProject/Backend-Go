package service

import (
	"errors"
	"greenenvironment/constant"
	"greenenvironment/features/admin"
	"greenenvironment/helper"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAdminRepository is a mock implementation of AdminRepositoryInterface
type MockAdminRepository struct {
	mock.Mock
}

func (m *MockAdminRepository) Login(adminData admin.Admin) (admin.Admin, error) {
	args := m.Called(adminData)
	return args.Get(0).(admin.Admin), args.Error(1)
}

func (m *MockAdminRepository) Update(adminData admin.AdminUpdate) (admin.Admin, error) {
	args := m.Called(adminData)
	return args.Get(0).(admin.Admin), args.Error(1)
}

func (m *MockAdminRepository) Delete(admin admin.Admin) error {
	args := m.Called(admin)
	return args.Error(0)
}

func (m *MockAdminRepository) GetAdminByID(id string) (admin.Admin, error) {
	args := m.Called(id)
	return args.Get(0).(admin.Admin), args.Error(1)
}

func (m *MockAdminRepository) IsEmailExist(email string) bool {
	args := m.Called(email)
	return args.Bool(0)
}

// MockJWT is a mock implementation of JWTInterface
type MockJWT struct {
	mock.Mock
}

func (m *MockJWT) GenerateUserJWT(userData helper.UserJWT) (string, error) {
	args := m.Called(userData)
	return args.String(0), args.Error(1)
}

func (m *MockJWT) GenerateAdminJWT(adminData helper.AdminJWT) (string, error) {
	args := m.Called(adminData)
	return args.String(0), args.Error(1)
}

func (m *MockJWT) GenerateUserToken(userData helper.UserJWT) string {
	args := m.Called(userData)
	return args.String(0)
}

func (m *MockJWT) GenerateAdminToken(adminData helper.AdminJWT) string {
	args := m.Called(adminData)
	return args.String(0)
}

func (m *MockJWT) ExtractUserToken(token *jwt.Token) map[string]interface{} {
	args := m.Called(token)
	return args.Get(0).(map[string]interface{})
}

func (m *MockJWT) ExtractAdminToken(token *jwt.Token) map[string]interface{} {
	args := m.Called(token)
	return args.Get(0).(map[string]interface{})
}

func (m *MockJWT) ValidateToken(token string) (*jwt.Token, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*jwt.Token), args.Error(1)
}

// setupTest creates new instances of mocks and service for testing
func setupTest() (*MockAdminRepository, *MockJWT, admin.AdminServiceInterface) {
	mockRepo := new(MockAdminRepository)
	mockJWT := new(MockJWT)
	adminService := NewAdminService(mockRepo, mockJWT)
	return mockRepo, mockJWT, adminService
}

func TestAdminService_Login(t *testing.T) {
	mockRepo, mockJWT, adminService := setupTest()

	t.Run("Success Login", func(t *testing.T) {
		inputAdmin := admin.Admin{
			Email:    "admin@test.com",
			Password: "password123",
		}

		expectedAdmin := admin.Admin{
			ID:       "1",
			Name:     "Admin Test",
			Email:    "admin@test.com",
			Username: "admintest",
		}

		expectedToken := "jwt.token.here"

		mockRepo.On("Login", inputAdmin).Return(expectedAdmin, nil).Once()
		mockJWT.On("GenerateAdminJWT", mock.AnythingOfType("helper.AdminJWT")).Return(expectedToken, nil).Once()

		result, err := adminService.Login(inputAdmin)

		assert.NoError(t, err)
		assert.Equal(t, expectedToken, result.Token)
		mockRepo.AssertExpectations(t)
		mockJWT.AssertExpectations(t)
	})

	t.Run("Failed Login - Repository Error", func(t *testing.T) {
		inputAdmin := admin.Admin{
			Email:    "admin@test.com",
			Password: "wrongpassword",
		}

		mockRepo.On("Login", inputAdmin).Return(admin.Admin{}, errors.New("invalid credentials")).Once()

		result, err := adminService.Login(inputAdmin)

		assert.Error(t, err)
		assert.Empty(t, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failed Login - JWT Generation Error", func(t *testing.T) {
		inputAdmin := admin.Admin{
			Email:    "admin@test.com",
			Password: "password123",
		}

		expectedAdmin := admin.Admin{
			ID:       "1",
			Name:     "Admin Test",
			Email:    "admin@test.com",
			Username: "admintest",
		}

		mockRepo.On("Login", inputAdmin).Return(expectedAdmin, nil).Once()
		mockJWT.On("GenerateAdminJWT", mock.AnythingOfType("helper.AdminJWT")).Return("", errors.New("jwt error")).Once()

		result, err := adminService.Login(inputAdmin)

		assert.Error(t, err)
		assert.Empty(t, result)
		mockRepo.AssertExpectations(t)
		mockJWT.AssertExpectations(t)
	})
}

func TestAdminService_Update(t *testing.T) {
	mockRepo, mockJWT, adminService := setupTest()

	t.Run("Success Update", func(t *testing.T) {
		inputAdmin := admin.AdminUpdate{
			ID:       "1",
			Name:     "Updated Admin",
			Email:    "updated@test.com",
			Username: "updatedadmin",
			Password: "newpassword123",
		}

		outputAdmin := admin.Admin{
			ID:        "1",
			Name:      "Updated Admin",
			Email:     "updated@test.com",
			Username:  "updatedadmin",
			Password:  "newpassword123",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		expectedToken := "new.jwt.token"

		mockRepo.On("Update", mock.AnythingOfType("admin.AdminUpdate")).Return(outputAdmin, nil).Once()
		mockJWT.On("GenerateAdminJWT", mock.AnythingOfType("helper.AdminJWT")).Return(expectedToken, nil).Once()

		result, err := adminService.Update(inputAdmin)

		assert.NoError(t, err)
		assert.Equal(t, expectedToken, result.Token)
		assert.Equal(t, inputAdmin.Name, result.Name)
		assert.Equal(t, inputAdmin.Email, result.Email)
		assert.Equal(t, inputAdmin.Username, result.Username)
		mockRepo.AssertExpectations(t)
		mockJWT.AssertExpectations(t)
	})

	t.Run("Invalid Username", func(t *testing.T) {
		inputAdmin := admin.AdminUpdate{
			Username: "i;lid",
		}

		result, err := adminService.Update(inputAdmin)

		assert.Error(t, err)
		assert.Equal(t, constant.ErrInvalidUsername, err)
		assert.Empty(t, result)
	})

	t.Run("Invalid Name", func(t *testing.T) {
		inputAdmin := admin.AdminUpdate{
			Username: "validuser",
			Name:     "",
		}

		result, err := adminService.Update(inputAdmin)

		assert.Error(t, err)
		assert.Equal(t, constant.ErrFieldData, err)
		assert.Empty(t, result)
	})
}

func TestAdminService_Delete(t *testing.T) {
	mockRepo, _, adminService := setupTest()

	t.Run("Success Delete", func(t *testing.T) {
		inputAdmin := admin.Admin{
			ID: "1",
		}

		mockRepo.On("Delete", inputAdmin).Return(nil).Once()

		err := adminService.Delete(inputAdmin)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Empty ID", func(t *testing.T) {
		inputAdmin := admin.Admin{
			ID: "",
		}

		err := adminService.Delete(inputAdmin)

		assert.Error(t, err)
		assert.Equal(t, constant.ErrDeleteUser, err)
	})

	t.Run("Repository Error", func(t *testing.T) {
		inputAdmin := admin.Admin{
			ID: "1",
		}

		mockRepo.On("Delete", inputAdmin).Return(errors.New("delete error")).Once()

		err := adminService.Delete(inputAdmin)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestAdminService_GetAdminData(t *testing.T) {
	mockRepo, _, adminService := setupTest()

	t.Run("Success Get Admin Data", func(t *testing.T) {
		inputAdmin := admin.Admin{
			ID: "1",
		}

		expectedAdmin := admin.Admin{
			ID:       "1",
			Name:     "Admin Test",
			Email:    "admin@test.com",
			Username: "admintest",
		}

		mockRepo.On("GetAdminByID", inputAdmin.ID).Return(expectedAdmin, nil).Once()

		result, err := adminService.GetAdminData(inputAdmin)

		assert.NoError(t, err)
		assert.Equal(t, expectedAdmin, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Not Found", func(t *testing.T) {
		inputAdmin := admin.Admin{
			ID: "999",
		}

		mockRepo.On("GetAdminByID", inputAdmin.ID).Return(admin.Admin{}, errors.New("admin not found")).Once()

		result, err := adminService.GetAdminData(inputAdmin)

		assert.Error(t, err)
		assert.Empty(t, result)
		mockRepo.AssertExpectations(t)
	})
}

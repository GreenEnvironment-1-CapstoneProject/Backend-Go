package helper

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/mock"
)

type MockJWTInterface struct {
	mock.Mock
}

func (m *MockJWTInterface) GenerateUserJWT(user UserJWT) (string, error) {
	args := m.Called(user)
	return args.String(0), args.Error(1)
}

func (m *MockJWTInterface) GenerateAdminJWT(user AdminJWT) (string, error) {
	args := m.Called(user)
	return args.String(0), args.Error(1)
}

func (m *MockJWTInterface) GenerateUserToken(user UserJWT) string {
	args := m.Called(user)
	return args.String(0)
}

func (m *MockJWTInterface) GenerateAdminToken(user AdminJWT) string {
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

package helper

import "github.com/stretchr/testify/mock"

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
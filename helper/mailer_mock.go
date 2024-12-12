package helper

import "github.com/stretchr/testify/mock"

type MockMailerInterface struct {
	mock.Mock
}

func (m *MockMailerInterface) Send(email, otpCode, subject string) error {
	args := m.Called(email, otpCode, subject)
	return args.Error(0)
}

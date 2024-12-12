package helper

import (
	"time"

	"github.com/stretchr/testify/mock"
)

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

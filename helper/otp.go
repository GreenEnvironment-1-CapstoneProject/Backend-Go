package helper

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type OTPInterface interface {
	GenerateOTP() string
	OTPExpiration(durationMinutes int) time.Time
}

type OTP struct {
	rng      *rand.Rand
	rngMutex sync.Mutex
}

func NewOTP() OTPInterface {
	return &OTP{
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (o *OTP) GenerateOTP() string {
	o.rngMutex.Lock()
	defer o.rngMutex.Unlock()
	return fmt.Sprintf("%06d", o.rng.Intn(1000000))
}

func (o *OTP) OTPExpiration(durationMinutes int) time.Time {
	return time.Now().Add(time.Duration(durationMinutes) * time.Minute)
}
package users

import (
	"time"

	"github.com/labstack/echo/v4"
)

type User struct {
	ID            string
	Username      string
	Password      string
	Name          string
	Email         string
	Address       string
	Gender        string
	Phone         string
	Exp           int
	Coin          int
	AvatarURL     string
	Is_Membership bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type UserLogin struct {
	Email    string
	Password string
	Token    string
}

type UserUpdate struct {
	ID        string
	Username  string
	Password  string
	Name      string
	Email     string
	Address   string
	Gender    string
	Phone     string
	AvatarURL string
	Token     string
}

type UpdateUserByAdmin struct {
	ID       string
	Name     string
	Address  string
	Gender   string
	Phone    string
	UpdateAt time.Time
}

type VerifyOTP struct {
	Email     string
	OTP       string
	ExpiredAt time.Time
}

type PasswordUpdate struct {
	Email       string
	OldPassword string
	NewPassword string
	OTP         string
}

type VerifyRegister struct {
	Name     string
	Email    string
	Password string
	OTP      string
}

type TemporaryUser struct {
	ID        string
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
}

type UserRepoInterface interface {
	Register(User) (User, error)
	SaveTemporaryUser(user TemporaryUser) error
	GetTemporaryUserByEmail(email string) (TemporaryUser, error)
	DeleteTemporaryUserByEmail(email string) error
	GetVerifyOTP(otp string) (VerifyOTP, error)
	DeleteVerifyOTP(otp string) error
	ValidateOTPByOTP(otp string) bool
	GetEmailByLatestOTP() (string, error)
	DeleteVerifyOTPByEmail(email string) error
	Login(User) (User, error)
	UpdateUserInfo(user UserUpdate) (User, error)
	Delete(User) error
	IsEmailExist(email string) bool
	GetUserByID(id string) (User, error)
	GetUserByEmail(email string) (User, error)
	UpdateAvatar(userID, avatarURL string) error

	SaveOTP(email, otp string, expiration time.Time) error
	ValidateOTP(email, otp string) bool
	UpdatePassword(email, hashedPassword string) error

	// Admin
	GetUserByIDForAdmin(id string) (User, error)
	GetAllUsersForAdmin() ([]User, error)
	UpdateUserForAdmin(UpdateUserByAdmin) error
	DeleteUserForAdmin(userID string) error
	GetAllByPageForAdmin(page int, limit int) ([]User, int, error)
}

type UserServiceInterface interface {
	RequestRegisterOTP(name, email, password string) error
	VerifyRegisterOTP(otp string) (User, error)
	IsEmailExist(email string) bool
	RequestPasswordResetOTP(email string) error
	VerifyPasswordResetOTP(otp string) error
	ResetPassword(newPassword string) error
	Login(User) (UserLogin, error)
	RegisterOrLoginGoogle(User) (User, error)
	UpdateUserInfo(user UserUpdate) error
	GetUserData(User) (User, error)
	Delete(User) error
	UpdateAvatar(userID, avatarURL string) error

	RequestPasswordUpdateOTP(email string) error
	UpdatePassword(update PasswordUpdate) error

	// Admin
	GetUserByIDForAdmin(id string) (User, error)
	UpdateUserForAdmin(UpdateUserByAdmin) error
	DeleteUserForAdmin(userID string) error
	GetAllByPageForAdmin(page int, limit int) ([]User, int, error)
}

type UserControllerInterface interface {
	RequestRegisterOTP(c echo.Context) error
	VerifyRegisterOTP(c echo.Context) error
	ForgotPasswordRequest(c echo.Context) error
	VerifyForgotPasswordOTP(c echo.Context) error
	ResetPassword(c echo.Context) error
	Login(c echo.Context) error
	GoogleLogin(c echo.Context) error
	GoogleCallback(c echo.Context) error
	UpdateUserInfo(c echo.Context) error
	GetUserData(c echo.Context) error
	Delete(c echo.Context) error
	UpdateAvatar(c echo.Context) error

	RequestPasswordUpdateOTP(c echo.Context) error
	UpdateUserPassword(c echo.Context) error

	// Admin
	GetAllUsersForAdmin(c echo.Context) error
	GetUserByIDForAdmin(c echo.Context) error
	UpdateUserForAdmin(c echo.Context) error
	DeleteUserForAdmin(c echo.Context) error
}

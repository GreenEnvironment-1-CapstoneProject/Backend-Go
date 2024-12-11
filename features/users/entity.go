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

type UserRepoInterface interface {
	Register(User) (User, error)
	Login(User) (User, error)
	UpdateUserInfo(user UserUpdate) (User, error)
	Delete(User) error
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
	Register(User) (User, error)
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
	Register(c echo.Context) error
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

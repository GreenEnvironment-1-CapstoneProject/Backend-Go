package controller

type UserRegisterRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"email,required"`
	Password string `json:"password" validate:"required"`
}

type UserVerifyRegisterRequest struct {
	OTP      string `json:"otp" validate:"required"`
}

type UserLoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserUpdateRequest struct {
	Name    string `json:"name" validate:"required"`
	Address string `json:"address" validate:"required"`
	Gender  string `json:"gender" validate:"required"`
	Phone   string `json:"phone" validate:"required"`
}

type UserPasswordUpdateRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required"`
	OTP         string `json:"otp" validate:"required"`
}

// Admin
type UserbyAdminRequest struct {
	Name    string `json:"name" validate:"required"`
	Address string `json:"address" validate:"required"`
	Gender  string `json:"gender" validate:"required"`
	Phone   string `json:"phone" validate:"required"`
}
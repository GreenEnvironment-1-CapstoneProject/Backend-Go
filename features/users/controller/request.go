package controller

type UserRegisterRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"email,required"`
	Password string `json:"password" validate:"required"`
}

type UserLoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserUpdateRequest struct {
	Name      string `json:"name" validate:"required"`
	Email     string `json:"email" validate:"required"`
	Username  string `json:"username" validate:"required"`
	Address   string `json:"address" validate:"required"`
	Gender    string `json:"gender" validate:"required"`
	Phone     string `json:"phone" validate:"required"`
	Password  string `json:"password" validate:"required"`
	AvatarURL string `json:"avatar_url" validate:"required"`
}

package controller

type UserRegisterResponse struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Address      string `json:"address"`
	Gender       string `json:"gender"`
	Phone        string `json:"phone"`
	Exp          int    `json:"exp"`
	Coin         int    `json:"coin"`
	AvatarURL    string `json:"avatar_url"`
	IsMembership bool   `json:"is_membership"`
}

type UserLoginResponse struct {
	Token string `json:"token"`
}

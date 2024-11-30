package controller

type UserRegisterResponse struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	Address       string `json:"address"`
	Gender        string `json:"gender"`
	Phone         string `json:"phone"`
	Exp           int    `json:"exp"`
	Coin          int    `json:"coin"`
	AvatarURL     string `json:"avatar_url"`
	Is_Membership bool   `json:"is_membership"`
}

type UserLoginResponse struct {
	Token string `json:"token"`
}

type UserInfoResponse struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Username      string `json:"username"`
	Email         string `json:"email"`
	Address       string `json:"address"`
	Gender        string `json:"gender"`
	Phone         string `json:"phone"`
	Coin          int    `json:"coin"`
	Exp           int    `json:"exp"`
	Is_Membership bool   `json:"is_membership"`
	AvatarURL     string `json:"avatar_url"`
}

type UserUpdateResponse struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Address  string `json:"address"`
	Gender   string `json:"gender"`
	Phone    string `json:"phone"`
}

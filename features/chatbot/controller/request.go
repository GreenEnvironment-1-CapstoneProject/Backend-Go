package controller

type ChatbotRequest struct {
	ID      string `json:"id"`
	Message string `json:"message" validate:"required"`
}

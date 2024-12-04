package chatbot

import (
	"time"

	"github.com/labstack/echo/v4"
)

type Chatbot struct {
	ID        string
	ChatID    string
	Role      string
	Message   string
	CreatedAt time.Time
}

type Message struct {
	Role    string
	Message string
}

type ChatbotControllerInterface interface {
	Create(c echo.Context) error
	GetByID(c echo.Context) error
}

type ChatbotServiceInterface interface {
	Create(chat Chatbot) (Chatbot, error)
	GetByID(id string) ([]Chatbot, error)
}

type ChatbotRepositoryInterface interface {
	Create(chat Chatbot) (Chatbot, error)
	GetByID(id string) ([]Chatbot, error)
}

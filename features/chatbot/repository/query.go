package repository

import (
	"greenenvironment/features/chatbot"

	"gorm.io/gorm"
)

type ChatbotRepository struct {
	DB *gorm.DB
}

func NewChatbotRepository(db *gorm.DB) chatbot.ChatbotRepositoryInterface {
	return &ChatbotRepository{
		DB: db,
	}
}

func (cr *ChatbotRepository) Create(chatbots chatbot.Chatbot) (chatbot.Chatbot, error) {
	result := cr.DB.Create(&chatbots)
	if result.Error != nil {
		return chatbot.Chatbot{}, result.Error
	}
	return chatbots, nil
}

func (cr *ChatbotRepository) GetByID(id string) ([]chatbot.Chatbot, error) {
	var chatbots []chatbot.Chatbot
	result := cr.DB.Where("chat_id = ?", id).
		Order("created_at ASC").
		Find(&chatbots)
	if result.Error != nil {
		return []chatbot.Chatbot{}, result.Error
	}
	return chatbots, nil
}

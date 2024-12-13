package service

import (
	"errors"
	"greenenvironment/features/chatbot"
	"testing"

	"github.com/google/uuid"
	openai "github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockChatbotRepository struct {
	mock.Mock
}

func (m *MockChatbotRepository) Create(chatbotDTO chatbot.Chatbot) (chatbot.Chatbot, error) {
	args := m.Called(chatbotDTO)
	return args.Get(0).(chatbot.Chatbot), args.Error(1)
}

func (m *MockChatbotRepository) GetByID(id string) ([]chatbot.Chatbot, error) {
	args := m.Called(id)
	return args.Get(0).([]chatbot.Chatbot), args.Error(1)
}

type MockOpenAIService struct {
	mock.Mock
}

func (m *MockOpenAIService) CreateChatCompletion(messages []openai.ChatCompletionMessage) (string, error) {
	args := m.Called(messages)
	return args.String(0), args.Error(1)
}

func TestCreateChatbot(t *testing.T) {
	mockRepo := new(MockChatbotRepository)
	mockOpenAI := new(MockOpenAIService)
	service := NewChatbotService(mockRepo, mockOpenAI)

	t.Run("success", func(t *testing.T) {
		chatID := uuid.New().String()
		chatbotIn := chatbot.Chatbot{
			ChatID:  chatID,
			Message: "Hello",
		}
		expectedResponse := chatbot.Chatbot{
			ID:      uuid.New().String(),
			ChatID:  chatID,
			Role:    "assistant",
			Message: "Hi there!",
		}

		mockRepo.On("GetByID", chatID).Return([]chatbot.Chatbot{}, nil)
		mockRepo.On("Create", mock.Anything).Return(chatbotIn, nil).Once()
		mockOpenAI.On("CreateChatCompletion", mock.Anything).Return("Hi there!", nil).Once()
		mockRepo.On("Create", mock.Anything).Return(expectedResponse, nil).Once()

		result, err := service.Create(chatbotIn)
		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, result)
		mockRepo.AssertExpectations(t)
		mockOpenAI.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		chatID := uuid.New().String()
		chatbotIn := chatbot.Chatbot{
			ChatID:  chatID,
			Message: "Hello",
		}

		mockRepo.On("GetByID", chatID).Return([]chatbot.Chatbot{}, nil)
		mockRepo.On("Create", mock.Anything).Return(chatbot.Chatbot{}, errors.New("repository error")).Once()

		result, err := service.Create(chatbotIn)
		assert.Error(t, err)
		assert.Equal(t, chatbot.Chatbot{}, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("openai error", func(t *testing.T) {
		chatID := uuid.New().String()
		chatbotIn := chatbot.Chatbot{
			ChatID:  chatID,
			Message: "Hello",
		}

		mockRepo.On("GetByID", chatID).Return([]chatbot.Chatbot{}, nil)
		mockRepo.On("Create", mock.Anything).Return(chatbotIn, nil).Once()
		mockOpenAI.On("CreateChatCompletion", mock.Anything).Return("", errors.New("openai error")).Once()

		result, err := service.Create(chatbotIn)
		assert.Error(t, err)
		assert.Equal(t, chatbot.Chatbot{}, result)
		mockRepo.AssertExpectations(t)
		mockOpenAI.AssertExpectations(t)
	})
}

func TestGetByID(t *testing.T) {
	mockRepo := new(MockChatbotRepository)
	service := NewChatbotService(mockRepo, nil)

	t.Run("success", func(t *testing.T) {
		chatID := uuid.New().String()
		expectedChats := []chatbot.Chatbot{
			{
				ID:      uuid.New().String(),
				ChatID:  chatID,
				Role:    "user",
				Message: "Hello",
			},
		}

		mockRepo.On("GetByID", chatID).Return(expectedChats, nil).Once()

		result, err := service.GetByID(chatID)
		assert.NoError(t, err)
		assert.Equal(t, expectedChats, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		chatID := uuid.New().String()

		mockRepo.On("GetByID", chatID).Return([]chatbot.Chatbot{}, errors.New("repository error")).Once()

		result, err := service.GetByID(chatID)
		assert.Error(t, err)
		assert.Equal(t, []chatbot.Chatbot{}, result)
		mockRepo.AssertExpectations(t)
	})
}

package service

import (
	"greenenvironment/features/chatbot"
	openaiService "greenenvironment/utils/openai"

	"github.com/google/uuid"
	openai "github.com/sashabaranov/go-openai"
)

type ChatbotService struct {
	chatbotRepo chatbot.ChatbotRepositoryInterface
	openAI      openaiService.OpenAIInterface
	systemRole  string
}

func NewChatbotService(d chatbot.ChatbotRepositoryInterface, openAI openaiService.OpenAIInterface) chatbot.ChatbotServiceInterface {
	return &ChatbotService{
		chatbotRepo: d,
		openAI:      openAI,
		systemRole: "Anda adalah seorang ahli dalam penghijauan lingkungan. " +
			"Anda akan diberikan pertanyaan tentang lingkungan dan harus menjawabnya " +
			"dengan relevansi tinggi terhadap lingkungan, produk eco-friendly, atau tantangan lingkungan hijau.",
	}
}

func (cs *ChatbotService) Create(chatbots chatbot.Chatbot) (chatbot.Chatbot, error) {
	payloadChatbot := []chatbot.Message{
		{
			Role:    "system",
			Message: cs.systemRole,
		},
	}

	var ID string

	if chatbots.ChatID == "" {
		ID = uuid.New().String()
	} else {
		ID = chatbots.ChatID
		previousChats, err := cs.chatbotRepo.GetByID(ID)
		if err != nil {
			return chatbot.Chatbot{}, err
		}
		for _, v := range previousChats {
			payloadChatbot = append(payloadChatbot, chatbot.Message{
				Role:    v.Role,
				Message: v.Message,
			})
		}
	}

	chatbots.ID = uuid.New().String()
	chatbots.ChatID = ID
	chatbots.Role = "user"

	if _, err := cs.chatbotRepo.Create(chatbots); err != nil {
		return chatbot.Chatbot{}, err
	}

	payloadChatbot = append(payloadChatbot, chatbot.Message{
		Role:    "user",
		Message: chatbots.Message,
	})

	openAIPayload := make([]openai.ChatCompletionMessage, len(payloadChatbot))
	for i, msg := range payloadChatbot {
		openAIPayload[i] = openai.ChatCompletionMessage{
			Role:    msg.Role,
			Content: msg.Message,
		}
	}

	assistantMessage, err := cs.openAI.CreateChatCompletion(openAIPayload)
	if err != nil {
		return chatbot.Chatbot{}, err
	}

	assistantResponse := chatbot.Chatbot{
		ID:      uuid.New().String(),
		ChatID:  ID,
		Role:    "assistant",
		Message: assistantMessage,
	}

	res, err := cs.chatbotRepo.Create(assistantResponse)
	if err != nil {
		return chatbot.Chatbot{}, err
	}

	return res, nil
}

func (cs *ChatbotService) GetByID(id string) ([]chatbot.Chatbot, error) {
	res, err := cs.chatbotRepo.GetByID(id)
	if err != nil {
		return []chatbot.Chatbot{}, err
	}
	return res, nil
}

package openai

import (
	"context"
	"errors"

	openai "github.com/sashabaranov/go-openai"
)

type OpenAIInterface interface {
	CreateChatCompletion(messages []openai.ChatCompletionMessage) (string, error)
}

type OpenAIService struct {
	client *openai.Client
	model  string
}

func NewOpenAIService(apiKey string) OpenAIInterface {
	client := openai.NewClient(apiKey)
	return &OpenAIService{
		client: client,
		model:  openai.GPT4,
	}
}

func (o *OpenAIService) CreateChatCompletion(messages []openai.ChatCompletionMessage) (string, error) {
	resp, err := o.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    o.model,
			Messages: messages,
		},
	)
	if err != nil {
		return "", err
	}
	if len(resp.Choices) == 0 {
		return "", errors.New("no response from OpenAI API")
	}

	return resp.Choices[0].Message.Content, nil
}

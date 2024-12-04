package controller

import (
	"greenenvironment/features/chatbot"
	"greenenvironment/helper"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ChatbotController struct {
	chatbotService chatbot.ChatbotServiceInterface
	jwtService     helper.JWTInterface
}

func NewChatbotController(s chatbot.ChatbotServiceInterface, j helper.JWTInterface) chatbot.ChatbotControllerInterface {
	return &ChatbotController{
		chatbotService: s,
		jwtService:     j,
	}
}

// Create Chatbot
// @Summary      Create a chatbot conversation
// @Description  Add a user message and get a response from the chatbot
// @Tags         Chatbot
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true   "Bearer Token"
// @Param        body           body      ChatbotRequest  true   "Chatbot message payload"
// @Success      201  {object}  helper.Response{data=ChatbotResponse} "Success Create Chatbot"
// @Failure      400  {object}  helper.Response{data=string} "Invalid request payload"
// @Failure      401  {object}  helper.Response{data=string} "Unauthorized"
// @Failure      500  {object}  helper.Response{data=string} "Internal server error"
// @Router       /chatbots [post]
func (cc *ChatbotController) Create(c echo.Context) error {
	tokenString := c.Request().Header.Get("Authorization")
	_, err := cc.jwtService.ValidateToken(tokenString)
	if err != nil {
		return helper.UnauthorizedError(c)
	}

	var chatBotRequest ChatbotRequest

	if err := c.Bind(&chatBotRequest); err != nil {
		return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "error bad request", nil))
	}

	if err := c.Validate(chatBotRequest); err != nil {
		return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, err.Error(), nil))
	}

	chatBot := chatbot.Chatbot{
		ChatID:  chatBotRequest.ID,
		Message: chatBotRequest.Message,
	}

	res, err := cc.chatbotService.Create(chatBot)
	if err != nil {
		return c.JSON(helper.ConvertResponseCode(err), helper.FormatResponse(false, err.Error(), nil))
	}

	var chatBotResponse ChatbotResponse
	chatBotResponse.ID = res.ID
	chatBotResponse.ChatID = res.ChatID
	chatBotResponse.Role = res.Role
	chatBotResponse.Message = res.Message
	chatBotResponse.CreatedAt = res.CreatedAt.Format("02/01/2006")

	return c.JSON(http.StatusCreated, helper.FormatResponse(true, "Success Create Chatbot", chatBotResponse))
}

// Get Chatbot by ID
// @Summary      Retrieve chatbot conversation by ChatID
// @Description  Fetch all messages in a chatbot session using the ChatID
// @Tags         Chatbot
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true   "Bearer Token"
// @Param        chatID         path      string  true   "Chat ID"
// @Success      200  {object}  helper.Response{data=[]ChatbotResponse} "Success Get Chatbot"
// @Failure      401  {object}  helper.Response{data=string} "Unauthorized"
// @Failure      404  {object}  helper.Response{data=string} "Chatbot not found"
// @Failure      500  {object}  helper.Response{data=string} "Internal server error"
// @Router       /chatbots/{chatID} [get]
func (cc *ChatbotController) GetByID(c echo.Context) error {
	chatID := c.Param("chatID")
	tokenString := c.Request().Header.Get("Authorization")
	_, err := cc.jwtService.ValidateToken(tokenString)
	if err != nil {
		return helper.UnauthorizedError(c)
	}

	res, err := cc.chatbotService.GetByID(chatID)
	if err != nil {
		return c.JSON(helper.ConvertResponseCode(err), helper.FormatResponse(false, err.Error(), nil))
	}
	var chatBotResponse []ChatbotResponse
	for _, v := range res {
		chatBotResponse = append(chatBotResponse, ChatbotResponse{
			ID:        v.ID,
			ChatID:    v.ChatID,
			Role:      v.Role,
			Message:   v.Message,
			CreatedAt: v.CreatedAt.Format("02/01/2006"),
		})
	}
	return c.JSON(http.StatusOK, helper.FormatResponse(true, "Success Get Chatbot", chatBotResponse))
}

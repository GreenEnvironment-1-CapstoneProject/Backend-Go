package controller

import (
	"greenenvironment/constant"
	"greenenvironment/features/forum"
	"greenenvironment/helper"
	"greenenvironment/utils/storages"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ForumController struct {
	forumService forum.ForumServiceInterface
	jwtService   helper.JWTInterface
	storage      storages.StorageInterface
}

func NewForumController(s forum.ForumServiceInterface, j helper.JWTInterface, storage storages.StorageInterface) forum.ForumControllerInterface {
	return &ForumController{
		forumService: s,
		jwtService:   j,
		storage:      storage,
	}
}

// @Summary      Get All Forums
// @Description  Retrieve all forums with pagination
// @Tags         Forum
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true   "Bearer Token"
// @Param        page           query     int     false  "Page number (default: 1)"
// @Success      200  {object}  helper.MetadataResponse "get all forum successfully"
// @Failure      400  {object}  helper.Response "Bad Request"
// @Failure      401  {object}  helper.Response "Unauthorized"
// @Failure      500  {object}  helper.Response "Internal Server Error"
// @Router       /forums [get]
func (h *ForumController) GetAllForum(c echo.Context) error {
	tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
	if tokenString == "" {
		helper.UnauthorizedError(c)
	}

	_, err := h.jwtService.ValidateToken(tokenString)
	if err != nil {
		helper.UnauthorizedError(c)
	}

	pageStr := c.QueryParam("page")
	page := 1
	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, constant.ErrPageInvalid.Error(), nil))
		}
	}

	var totalPages int
	var forums []forum.Forum
	forums, totalPages, err = h.forumService.GetAllByPage(page)
	metadata := MetadataResponse{
		CurrentPage: page,
		TotalPage:   totalPages,
	}

	if err != nil {
		code, message := helper.HandleEchoError(err)
		return c.JSON(code, helper.FormatResponse(false, message, nil))
	}

	var response []ForumGetAllResponse
	for _, f := range forums {
		response = append(response, ForumGetAllResponse{
			ID:          f.ID,
			Title:       f.Title,
			Description: f.Description,
			View:        f.View,
			TopicImage:  f.TopicImage,
			CreatedAt:   f.CreatedAt.Format("02/01/2006"),
			UpdatedAt:   f.UpdatedAt.Format("02/01/2006"),
			Author:      Author{ID: f.User.ID, Name: f.User.Name, Username: f.User.Username, Email: f.User.Email, AvatarURL: f.User.AvatarURL},
		})
	}
	return c.JSON(http.StatusOK, helper.MetadataFormatResponse(true, "get all forum successfully", metadata, response))
}

// @Summary      Create Forum
// @Description  Create a new forum post with title, description, and optional topic image
// @Tags         Forum
// @Accept       multipart/form-data
// @Produce      json
// @Param        Authorization  header    string                  true   "Bearer Token"
// @Param        title          formData  string                  true   "Forum Title"
// @Param        description    formData  string                  true   "Forum Description"
// @Param        topic_image    formData  file                    false  "Forum Topic Image"
// @Success      201  {object}  helper.Response                  "create forum successfully"
// @Failure      400  {object}  helper.Response                  "Bad Request"
// @Failure      401  {object}  helper.Response                  "Unauthorized"
// @Failure      500  {object}  helper.Response                  "Internal Server Error"
// @Router       /forums [post]
func (h *ForumController) PostForum(c echo.Context) error {
	tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
	if tokenString == "" {
		helper.UnauthorizedError(c)
	}

	token, err := h.jwtService.ValidateToken(tokenString)
	if err != nil {
		helper.UnauthorizedError(c)
	}

	userData := h.jwtService.ExtractUserToken(token)
	userId := userData[constant.JWT_ID]
	var forumData CreateForumRequest
	c.Bind(&forumData)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, err.Error(), nil))
	}

	if err := c.Validate(forumData); err != nil {
		return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, err.Error(), nil))
	}

	forumData.Title = strings.TrimSpace(forumData.Title)
	forumData.Description = strings.TrimSpace(forumData.Description)

	file, err := c.FormFile("topic_image")
	var topic_image string
	if file != nil {

		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "topic_image file is required", nil))
		}

		src, err := h.storage.ImageValidation(file)
		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, err.Error(), nil))
		}

		topic_image, err = h.storage.UploadImageToCloudinary(src, "ecomate/topic_image/")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, "failed to upload avatar", nil))
		}
	}

	newForum := forum.Forum{
		ID:          uuid.New().String(),
		Title:       forumData.Title,
		Description: forumData.Description,
		UserID:      userId.(string),
		TopicImage:  topic_image,
	}
	err = h.forumService.PostForum(newForum)
	if err != nil {
		return c.JSON(helper.ConvertResponseCode(err), helper.FormatResponse(false, err.Error(), nil))
	}
	return c.JSON(http.StatusCreated, helper.FormatResponse(true, "create forum successfully", nil))
}

// @Summary      Get Forum by ID
// @Description  Retrieve detailed information about a forum, including messages
// @Tags         Forum
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string                  true   "Bearer Token"
// @Param        id             path      string                  true   "Forum ID"
// @Success      200  {object}  helper.Response{data=ForumGetDetailResponse} "get forum successfully"
// @Failure      401  {object}  helper.Response                  "Unauthorized"
// @Failure      404  {object}  helper.Response                  "Forum not found"
// @Router       /forums/{id} [get]
func (h *ForumController) GetForumByID(c echo.Context) error {
	tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
	if tokenString == "" {
		return helper.UnauthorizedError(c)
	}

	_, err := h.jwtService.ValidateToken(tokenString)
	if err != nil {
		return helper.UnauthorizedError(c)
	}

	forumID := c.Param("id")

	forum, err := h.forumService.GetForumByID(forumID)
	if err != nil {
		return c.JSON(http.StatusNotFound, helper.ObjectFormatResponse(false, err.Error(), nil))
	}

	messages, err := h.forumService.GetMessagesByForumID(forumID)
	if err != nil {
		return c.JSON(http.StatusNotFound, helper.FormatResponse(false, err.Error(), nil))
	}

	var messageResponses []MessageResponse
	for _, msg := range messages {
		messageResponses = append(messageResponses, MessageResponse{
			ID: msg.ID,
			User: AuthorMessage{
				ID:        msg.User.ID,
				Name:      msg.User.Name,
				Username:  msg.User.Username,
				Email:     msg.User.Email,
				AvatarURL: msg.User.AvatarURL,
			},
			Message:      msg.Message,
			MessageImage: msg.MessageImage,
			CreatedAt:    msg.CreatedAt.Format("02/01/2006"),
			UpdatedAt:    msg.UpdatedAt.Format("02/01/2006"),
		})
	}

	forumResponse := ForumGetDetailResponse{
		ID:            forum.ID,
		Title:         forum.Title,
		Description:   forum.Description,
		TopicImage:    forum.TopicImage,
		View:          forum.View,
		Author:        Author{ID: forum.User.ID, Name: forum.User.Name, Username: forum.User.Username, Email: forum.User.Email, AvatarURL: forum.User.AvatarURL},
		ForumMessages: messageResponses,
		CreatedAt:     forum.CreatedAt.Format("02/01/2006"),
		UpdatedAt:     forum.UpdatedAt.Format("02/01/2006"),
	}

	return c.JSON(http.StatusOK, helper.ObjectFormatResponse(true, "get forum successfully", forumResponse))
}

// @Summary      Update Forum
// @Description  Update the details of an existing forum, including title, description, and topic image.
// @Tags         Forum
// @Accept       multipart/form-data
// @Produce      json
// @Param        Authorization  header    string  true   "Bearer Token"
// @Param        id             path      string  true   "Forum ID"
// @Param        title          formData  string  true   "New title of the forum"
// @Param        description    formData  string  true   "New description of the forum"
// @Param        topic_image    formData  file    false  "New topic image for the forum"
// @Success      200  {object}  helper.Response  "update forum successfully"
// @Failure      400  {object}  helper.Response             "Bad Request"
// @Failure      401  {object}  helper.Response             "Unauthorized"
// @Failure      403  {object}  helper.Response             "Forbidden"
// @Failure      404  {object}  helper.Response             "Forum not found"
// @Failure      500  {object}  helper.Response             "Internal Server Error"
// @Router       /forums/{id} [put]
func (h *ForumController) UpdateForum(c echo.Context) error {
	tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
	if tokenString == "" {
		return helper.UnauthorizedError(c)
	}

	token, err := h.jwtService.ValidateToken(tokenString)
	if err != nil {
		return helper.UnauthorizedError(c)
	}

	userData := h.jwtService.ExtractUserToken(token)
	userId := userData[constant.JWT_ID].(string)

	forumid := c.Param("id")
	existingForum, err := h.forumService.GetForumByID(forumid)
	if err != nil {
		return c.JSON(http.StatusNotFound, helper.FormatResponse(false, string(err.Error()), nil))
	}

	if existingForum.UserID != userId {
		return c.JSON(http.StatusForbidden, helper.FormatResponse(false, "error forbidden", nil))
	}
	var forumData EditForumRequest
	if err := c.Bind(&forumData); err != nil {
		code, message := helper.HandleEchoError(err)
		return c.JSON(code, helper.FormatResponse(false, message, nil))
	}

	if err := c.Validate(forumData); err != nil {
		return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, err.Error(), nil))
	}
	file, err := c.FormFile("topic_image")
	var topic_image string
	if file != nil {

		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "topic_image file is required", nil))
		}

		src, err := h.storage.ImageValidation(file)
		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, err.Error(), nil))
		}

		topic_image, err = h.storage.UploadImageToCloudinary(src, "ecomate/topic_image/")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, "failed to upload avatar", nil))
		}
	}

	forumData.Title = strings.TrimSpace(forumData.Title)
	forumData.Description = strings.TrimSpace(forumData.Description)
	forumsResponse := forum.EditForum{
		ID:          forumid,
		Title:       forumData.Title,
		Description: forumData.Description,
		TopicImage:  topic_image,
	}
	if err := h.forumService.UpdateForum(forumsResponse); err != nil {
		return c.JSON(helper.ConvertResponseCode(err), helper.FormatResponse(false, err.Error(), nil))
	}

	return c.JSON(http.StatusOK, helper.ObjectFormatResponse(true, "update forum successfully", nil))
}

// @Summary      Delete Forum
// @Description  Delete a forum by its ID. This action is restricted to users with admin role.
// @Tags         Forum
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true   "Bearer Token"
// @Param        id             path      string  true   "Forum ID"
// @Success      200  {object}  helper.Response  "delete forum successfully"
// @Failure      400  {object}  helper.Response  "Bad Request"
// @Failure      401  {object}  helper.Response  "Unauthorized"
// @Failure      403  {object}  helper.Response  "Forbidden"
// @Failure      404  {object}  helper.Response  "Forum not found"
// @Failure      500  {object}  helper.Response  "Internal Server Error"
// @Router       /forums/{id} [delete]
func (h *ForumController) DeleteForum(c echo.Context) error {
	tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
	if tokenString == "" {
		return helper.UnauthorizedError(c)
	}

	token, err := h.jwtService.ValidateToken(tokenString)
	if err != nil {
		return helper.UnauthorizedError(c)
	}

	userData := h.jwtService.ExtractUserToken(token)
	userId := userData[constant.JWT_ID].(string)

	forumID := c.Param("id")

	existingForum, err := h.forumService.GetForumByID(forumID)
	if err != nil {
		return c.JSON(http.StatusNotFound, helper.FormatResponse(false, string(err.Error()), nil))
	}

	if existingForum.UserID != userId {
		return c.JSON(http.StatusForbidden, helper.FormatResponse(false, "error forbidden", nil))
	}

	err = h.forumService.DeleteForum(forumID)
	if err != nil {
		return c.JSON(helper.ConvertResponseCode(err), helper.FormatResponse(false, err.Error(), nil))
	}

	return c.JSON(http.StatusOK, helper.FormatResponse(true, "delete forum successfully", nil))
}

// @Summary      Get Forums by User ID
// @Description  Retrieve a list of forums created by the authenticated user, with pagination support.
// @Tags         Forum
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true   "Bearer Token"
// @Param        page           query     int     false  "Page number for pagination (default is 1)"
// @Success      200  {object}  helper.MetadataResponse[data=[]ForumGetAllResponse]  "get forums by user id successfully"
// @Failure      400  {object}  helper.Response                                "Bad Request"
// @Failure      401  {object}  helper.Response                                "Unauthorized"
// @Failure      500  {object}  helper.Response                                "Internal Server Error"
// @Router       /forums/user [get]
func (h *ForumController) GetForumByUserID(c echo.Context) error {
	tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
	if tokenString == "" {
		helper.UnauthorizedError(c)
	}

	token, err := h.jwtService.ValidateToken(tokenString)
	if err != nil {
		return helper.UnauthorizedError(c)
	}

	userData := h.jwtService.ExtractUserToken(token)
	userId := userData[constant.JWT_ID].(string)

	pageStr := c.QueryParam("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}
	var totalPages int
	forums, totalPages, err := h.forumService.GetForumByUserID(userId, page)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, err.Error(), nil))
	}

	metadata := MetadataResponse{
		CurrentPage: page,
		TotalPage:   totalPages,
	}

	var forumsResponse []ForumGetAllResponse

	for _, forum := range forums {
		forumsResponse = append(forumsResponse, ForumGetAllResponse{
			ID:          forum.ID,
			Title:       forum.Title,
			Description: forum.Description,
			TopicImage:  forum.TopicImage,
			View:        forum.View,
			CreatedAt:   forum.CreatedAt.Format("02/01/2006"),
			UpdatedAt:   forum.UpdatedAt.Format("02/01/2006"),
			Author:      Author{ID: forum.User.ID, Name: forum.User.Name, Username: forum.User.Username, Email: forum.User.Email, AvatarURL: forum.User.AvatarURL},
		})
	}

	return c.JSON(http.StatusOK, helper.MetadataFormatResponse(true, "get forum by id successfully", metadata, forumsResponse))
}

// @Summary      Post Message to Forum
// @Description  Post a message to a specific forum. Optionally includes an image with the message.
// @Tags         Forum
// @Accept       multipart/form-data
// @Produce      json
// @Param        Authorization   header    string  true   "Bearer Token"
// @Param        ForumID         formData  string  true   "ID of the forum to post the message to"
// @Param        Messages        formData  string  true   "Content of the message"
// @Param        message_image   formData  file    false  "Optional image to include with the message"
// @Success      201  {object}   helper.Response  "create message successfully"
// @Failure      400  {object}   helper.Response  "Bad Request"
// @Failure      401  {object}   helper.Response  "Unauthorized"
// @Failure      404  {object}   helper.Response  "Forum not found"
// @Failure      500  {object}   helper.Response  "Internal Server Error"
// @Router       /forums/message [post]
func (h *ForumController) PostMessageForum(c echo.Context) error {
	tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
	if tokenString == "" {
		helper.UnauthorizedError(c)
	}

	token, err := h.jwtService.ValidateToken(tokenString)
	if err != nil {
		helper.UnauthorizedError(c)
	}

	userData := h.jwtService.ExtractUserToken(token)
	userId := userData[constant.JWT_ID].(string)

	var forumMessage CreateMessageForumRequest
	if err := c.Bind(&forumMessage); err != nil {
		code, message := helper.HandleEchoError(err)
		return c.JSON(code, helper.FormatResponse(false, message, nil))
	}
	if err := c.Validate(forumMessage); err != nil {
		return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, err.Error(), nil))
	}

	file, err := c.FormFile("message_image")
	var message_image string
	if file != nil {

		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "message_image file is required", nil))
		}

		src, err := h.storage.ImageValidation(file)
		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, err.Error(), nil))
		}

		message_image, err = h.storage.UploadImageToCloudinary(src, "ecomate/message_image/")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, "failed to upload avatar", nil))
		}
	}

	forumData := forum.MessageForum{
		ID:           uuid.New().String(),
		ForumID:      forumMessage.ForumID,
		UserID:       userId,
		Message:      forumMessage.Messages,
		MessageImage: message_image,
	}
	_, err = h.forumService.GetForumByID(forumData.ForumID)
	if err != nil {
		return c.JSON(http.StatusNotFound, helper.FormatResponse(false, err.Error(), nil))
	}

	err = h.forumService.PostMessageForum(forumData)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, err.Error(), nil))
	}

	return c.JSON(http.StatusCreated, helper.FormatResponse(true, "create message successfully", nil))
}

// @Summary      Delete Forum Message
// @Description  Delete a specific message in a forum. Only the message owner or an admin can delete it.
// @Tags         Forum
// @Accept       json
// @Produce      json
// @Param        Authorization   header    string  true   "Bearer Token"
// @Param        id              path      string  true   "ID of the forum message to delete"
// @Success      200  {object}   helper.Response  "delete message successfully"
// @Failure      401  {object}   helper.Response  "Unauthorized"
// @Failure      403  {object}   helper.Response  "Forbidden"
// @Failure      404  {object}   helper.Response  "Message not found"
// @Failure      500  {object}   helper.Response  "Internal Server Error"
// @Router       /forums/message/{id} [delete]
func (h *ForumController) DeleteMessageForum(c echo.Context) error {
	tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
	if tokenString == "" {
		helper.UnauthorizedError(c)
	}

	token, err := h.jwtService.ValidateToken(tokenString)
	if err != nil {
		helper.UnauthorizedError(c)
	}

	userData := h.jwtService.ExtractUserToken(token)
	userId := userData[constant.JWT_ID].(string)
	if userId == "" {
		return c.JSON(helper.ConvertResponseCode(err), helper.FormatResponse(false, constant.Unauthorized, nil))
	}
	isAdmin := userData[constant.JWT_ROLE] == constant.RoleAdmin

	messageForumID := c.Param("id")
	existingMessageForum, err := h.forumService.GetMessageForumByID(messageForumID)
	if err != nil {
		return c.JSON(http.StatusNotFound, helper.FormatResponse(false, err.Error(), nil))
	}

	if existingMessageForum.UserID != userId && !isAdmin {
		return c.JSON(http.StatusForbidden, helper.FormatResponse(false, constant.Unauthorized, nil))
	}

	err = h.forumService.DeleteMessageForum(messageForumID)
	if err != nil {
		return c.JSON(helper.ConvertResponseCode(err), helper.FormatResponse(false, err.Error(), nil))
	}

	return c.JSON(http.StatusOK, helper.FormatResponse(true, "delete message successfully", nil))
}

// @Summary      Update Forum Message
// @Description  Update the content or image of a specific message in a forum.
// @Tags         Forum
// @Accept       multipart/form-data
// @Produce      json
// @Param        Authorization   header    string  true   "Bearer Token"
// @Param        id              path      string  true   "ID of the message to update"
// @Param        Messages        formData  string  true   "Updated content of the message"
// @Param        message_image   formData  file    false  "Optional new image for the message"
// @Success      200  {object}   helper.Response					  "update message successfully"
// @Failure      400  {object}   helper.Response             "Bad Request"
// @Failure      401  {object}   helper.Response             "Unauthorized"
// @Failure      403  {object}   helper.Response             "Forbidden"
// @Failure      404  {object}   helper.Response             "Message not found"
// @Failure      500  {object}   helper.Response             "Internal Server Error"
// @Router       /forums/message/{id} [put]
func (h *ForumController) UpdateMessageForum(c echo.Context) error {
	tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
	if tokenString == "" {
		helper.UnauthorizedError(c)
	}

	token, err := h.jwtService.ValidateToken(tokenString)
	if err != nil {
		return helper.UnauthorizedError(c)
	}

	userData := h.jwtService.ExtractUserToken(token)
	userId := userData[constant.JWT_ID].(string)

	messageId := c.Param("id")
	existingMessage, err := h.forumService.GetMessageForumByID(messageId)
	if err != nil {
		return c.JSON(http.StatusNotFound, helper.FormatResponse(false, err.Error(), nil))
	}

	if existingMessage.UserID != userId {
		return c.JSON(http.StatusForbidden, helper.FormatResponse(false, "error forbidden", nil))
	}
	var messageForum EditMessageRequest
	if err := c.Bind(&messageForum); err != nil {
		code, message := helper.HandleEchoError(err)
		return c.JSON(code, helper.FormatResponse(false, message, nil))
	}

	if err := c.Validate(messageForum); err != nil {
		return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, err.Error(), nil))
	}
	file, err := c.FormFile("message_image")
	var message_image string
	if file != nil {

		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "message_image file is required", nil))
		}

		src, err := h.storage.ImageValidation(file)
		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, err.Error(), nil))
		}

		message_image, err = h.storage.UploadImageToCloudinary(src, "ecomate/message_image/")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, "failed to upload avatar", nil))
		}
	}

	messageResponse := forum.EditMessage{
		ID:           messageId,
		Message:      messageForum.Messages,
		MessageImage: message_image,
	}
	if err := h.forumService.UpdateMessageForum(messageResponse); err != nil {
		return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, err.Error(), nil))
	}

	return c.JSON(http.StatusOK, helper.ObjectFormatResponse(true, "update message successfully", nil))
}

// @Summary      Get Forum Message by ID
// @Description  Retrieve detailed information about a specific message in a forum.
// @Tags         Forum
// @Accept       json
// @Produce      json
// @Param        Authorization   header    string  true   "Bearer Token"
// @Param        id              path      string  true   "ID of the forum message to retrieve"
// @Success      200  {object}   helper.Response{data=MessageResponse}  "get message forum id successfully"
// @Failure      401  {object}   helper.Response                         "Unauthorized"
// @Failure      404  {object}   helper.Response                         "Message not found"
// @Failure      500  {object}   helper.Response                         "Internal Server Error"
// @Router       /forums/message/{id} [get]
func (h *ForumController) GetMessageForumByID(c echo.Context) error {
	tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
	if tokenString == "" {
		helper.UnauthorizedError(c)
	}

	_, err := h.jwtService.ValidateToken(tokenString)
	if err != nil {
		helper.UnauthorizedError(c)
	}

	messageId := c.Param("id")
	message, err := h.forumService.GetMessageForumByID(messageId)
	if err != nil {
		return c.JSON(http.StatusNotFound, helper.ObjectFormatResponse(false, "get message not found", nil))
	}

	messagesResponse := MessageResponse{
		ID:           message.ID,
		User:         AuthorMessage{ID: message.UserID, Name: message.User.Name, Username: message.User.Username, Email: message.User.Email, AvatarURL: message.User.AvatarURL},
		Message:      message.Message,
		MessageImage: message.MessageImage,
		CreatedAt:    message.CreatedAt.Format("02/01/2006"),
		UpdatedAt:    message.UpdatedAt.Format("02/01/2006"),
	}

	return c.JSON(http.StatusOK, helper.ObjectFormatResponse(true, "get message forum id successfully", messagesResponse))
}

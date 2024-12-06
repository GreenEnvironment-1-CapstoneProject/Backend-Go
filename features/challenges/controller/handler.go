package controller

import (
	"errors"
	"greenenvironment/constant"
	"greenenvironment/features/challenges"
	"greenenvironment/helper"
	"greenenvironment/utils/storages"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type ChallengeHandler struct {
	challengeService challenges.ChallengeServiceInterface
	jwt              helper.JWTInterface
	storage          storages.StorageInterface
}

func NewChallengeController(s challenges.ChallengeServiceInterface, j helper.JWTInterface, sto storages.StorageInterface) challenges.ChallengeControllerInterface {
	return &ChallengeHandler{
		challengeService: s,
		jwt:              j,
		storage:          sto,
	}
}

// Challenge
// @Summary      Create a new challenge
// @Description  Create a new challenge. Requires admin role.
// @Tags         Challenges (Admin)
// @Accept       multipart/form-data
// @Produce      json
// @Param        Authorization    header    string               true   "Bearer Token"
// @Param        challenge_img    formData  file                 true   "Challenge image file"
// @Param        title            formData  string               true   "Title of the challenge"
// @Param        difficulty       formData  string               true   "Difficulty level"
// @Param        description      formData  string               true   "Description of the challenge"
// @Param        duration_days    formData  int                  true   "Duration of the challenge in days"
// @Param        exp              formData  int                  true   "Experience points awarded for completing the challenge"
// @Param        coin             formData  int                  true   "Coins awarded for completing the challenge"
// @Param        impact_categories formData []string                true   "List of impact category IDs"
// @Success      201  {object}    helper.Response{data=string}   "Challenge created successfully"
// @Failure      400  {object}    helper.Response{data=string}   "Bad request"
// @Failure      401  {object}    helper.Response{data=string}   "Unauthorized"
// @Failure      500  {object}    helper.Response{data=string}   "Internal server error"
// @Router       /admin/challenges [post]
func (h *ChallengeHandler) Create(c echo.Context) error {
	tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
	if tokenString == "" {
		helper.UnauthorizedError(c)
	}

	token, err := h.jwt.ValidateToken(tokenString)
	if err != nil {
		helper.UnauthorizedError(c)
	}

	adminData := h.jwt.ExtractUserToken(token)
	adminRole := adminData[constant.JWT_ROLE].(string)
	if adminRole != constant.RoleAdmin {
		helper.UnauthorizedError(c)
	}
	adminId := adminData[constant.JWT_ID].(string)

	file, err := c.FormFile("challenge_img")
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "Challenge image file is required", nil))
	}

	src, err := h.storage.ImageValidation(file)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, err.Error(), nil))
	}

	challengeImgURL, err := h.storage.UploadImageToCloudinary(src, "ecomate/challenges/images/")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, "Failed to upload challenge image", nil))
	}

	var challengeRequest ChallengeRequest
	if err := c.Bind(&challengeRequest); err != nil {
		return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "Error bad request", nil))
	}

	if err := c.Validate(challengeRequest); err != nil {
		return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, err.Error(), nil))
	}

	challenge := challenges.Challenge{
		Author:       adminId,
		Title:        challengeRequest.Title,
		Difficulty:   challengeRequest.Difficulty,
		ChallengeImg: challengeImgURL,
		Description:  challengeRequest.Description,
		DurationDays: challengeRequest.DurationDays,
		Exp:          challengeRequest.Exp,
		Coin:         challengeRequest.Coin,
	}

	for _, category := range challengeRequest.ImpactCategories {
		challenge.ImpactCategories = append(challenge.ImpactCategories, challenges.ChallengeImpactCategory{
			ImpactCategoryID: category,
		})
	}

	err = h.challengeService.Create(challenge)
	if err != nil {
		return c.JSON(helper.ConvertResponseCode(err), helper.FormatResponse(false, err.Error(), nil))
	}

	return c.JSON(http.StatusCreated, helper.FormatResponse(true, "Success", nil))
}

// Get All Challenges
// @Summary      Get all challenges
// @Description  Retrieve all challenges with pagination.
// @Tags         Challenges (Admin)
// @Accept       json
// @Produce      json
// @Param        pages          query     int     false  "Page number"
// @Success      200  {object}  helper.MetadataResponse{data=[]ChallengeResponse} "Challenges retrieved successfully"
// @Failure      500  {object}  helper.Response{data=string} "Internal server error"
// @Router       /admin/challenges [get]
func (h *ChallengeHandler) GetAll(c echo.Context) error {
	tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
	if tokenString == "" {
		return helper.UnauthorizedError(c)
	}

	token, err := h.jwt.ValidateToken(tokenString)
	if err != nil {
		return helper.UnauthorizedError(c)
	}

	adminData := h.jwt.ExtractUserToken(token)
	adminRole := adminData[constant.JWT_ROLE].(string)
	if adminRole != constant.RoleAdmin {
		return helper.UnauthorizedError(c)
	}

	page, err := strconv.Atoi(c.QueryParam("pages"))
	if err != nil || page < 1 {
		page = 1
	}

	challenges, totalPages, err := h.challengeService.GetAllByPage(page)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, err.Error(), nil))
	}

	var response []interface{}
	for _, challenge := range challenges {
		response = append(response, new(ChallengeResponse).ToResponse(challenge))
	}

	metadata := map[string]interface{}{
		"TotalPage": totalPages,
		"Page":      page,
	}

	return c.JSON(http.StatusOK, helper.MetadataFormatResponse(true, "get all challenges successfully", metadata, response))
}

// Get Challenge by ID
// @Summary      Get a challenge by ID
// @Description  Retrieve a specific challenge by its unique identifier.
// @Tags         Challenges (Admin)
// @Accept       json
// @Produce      json
// @Param        id             path      string  true   "Challenge ID"
// @Success      200  {object}  helper.Response{data=ChallengeResponse} "Challenge retrieved successfully"
// @Failure      404  {object}  helper.Response{data=string} "Challenge not found"
// @Failure      500  {object}  helper.Response{data=string} "Internal server error"
// @Router       /admin/challenges/{id} [get]
func (h *ChallengeHandler) GetByID(c echo.Context) error {
	tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
	if tokenString == "" {
		return helper.UnauthorizedError(c)
	}

	token, err := h.jwt.ValidateToken(tokenString)
	if err != nil {
		return helper.UnauthorizedError(c)
	}

	adminData := h.jwt.ExtractUserToken(token)
	adminRole := adminData[constant.JWT_ROLE].(string)
	if adminRole != constant.RoleAdmin {
		return helper.UnauthorizedError(c)
	}

	id := c.Param("id")

	challenge, err := h.challengeService.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, helper.FormatResponse(false, err.Error(), nil))
	}

	response := new(ChallengeResponse).ToResponse(challenge)

	return c.JSON(http.StatusOK, helper.FormatResponse(true, "get challenge successfully", response))
}

// Update Challenge
// @Summary      Update a challenge
// @Description  Update challenge details, including impact categories. Requires admin role.
// @Tags         Challenges (Admin)
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true   "Bearer Token"
// @Param        id             path      string  true   "Challenge ID"
// @Param        body           body      ChallengeRequest  true   "Updated challenge data"
// @Success      200  {object}  helper.Response{data=string} "Challenge updated successfully"
// @Failure      400  {object}  helper.Response{data=string} "Bad request"
// @Failure      401  {object}  helper.Response{data=string} "Unauthorized"
// @Failure      404  {object}  helper.Response{data=string} "Challenge not found"
// @Failure      500  {object}  helper.Response{data=string} "Internal server error"
// @Router       /admin/challenges/{id} [put]
func (h *ChallengeHandler) Update(c echo.Context) error {
	tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
	if tokenString == "" {
		return helper.UnauthorizedError(c)
	}

	token, err := h.jwt.ValidateToken(tokenString)
	if err != nil {
		return helper.UnauthorizedError(c)
	}

	adminData := h.jwt.ExtractUserToken(token)
	adminRole := adminData[constant.JWT_ROLE].(string)
	if adminRole != constant.RoleAdmin {
		return helper.UnauthorizedError(c)
	}

	file, err := c.FormFile("challenge_img")
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "Challenge image file is required", nil))
	}

	src, err := h.storage.ImageValidation(file)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, err.Error(), nil))
	}

	challengeImgURL, err := h.storage.UploadImageToCloudinary(src, "ecomate/challenges/images/")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, "Failed to upload challenge image", nil))
	}

	challengeID := c.Param("id")
	var challengeRequest ChallengeRequest

	if err := c.Bind(&challengeRequest); err != nil {
		return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "Invalid request", nil))
	}

	if err := c.Validate(challengeRequest); err != nil {
		return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, err.Error(), nil))
	}

	challenge := challenges.Challenge{
		ID:           challengeID,
		Title:        challengeRequest.Title,
		Difficulty:   challengeRequest.Difficulty,
		ChallengeImg: challengeImgURL,
		Description:  challengeRequest.Description,
		DurationDays: challengeRequest.DurationDays,
		Exp:          challengeRequest.Exp,
		Coin:         challengeRequest.Coin,
	}

	for _, category := range challengeRequest.ImpactCategories {
		challenge.ImpactCategories = append(challenge.ImpactCategories, challenges.ChallengeImpactCategory{
			ImpactCategoryID: category,
		})
	}

	err = h.challengeService.Update(challenge)
	if err != nil {
		return c.JSON(helper.ConvertResponseCode(err), helper.FormatResponse(false, err.Error(), nil))
	}

	return c.JSON(http.StatusOK, helper.FormatResponse(true, "Challenge updated successfully", nil))
}

// Delete Challenge
// @Summary      Delete a challenge
// @Description  Remove a challenge by its ID. Requires admin role.
// @Tags         Challenges (Admin)
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true   "Bearer Token"
// @Param        id             path      string  true   "Challenge ID"
// @Success      200  {object}  helper.Response{data=string} "Challenge deleted successfully"
// @Failure      401  {object}  helper.Response{data=string} "Unauthorized"
// @Failure      404  {object}  helper.Response{data=string} "Challenge not found"
// @Failure      500  {object}  helper.Response{data=string} "Internal server error"
// @Router       /admin/challenges/{id} [delete]
func (h *ChallengeHandler) Delete(c echo.Context) error {
	tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
	if tokenString == "" {
		return helper.UnauthorizedError(c)
	}

	token, err := h.jwt.ValidateToken(tokenString)
	if err != nil {
		return helper.UnauthorizedError(c)
	}

	adminData := h.jwt.ExtractUserToken(token)
	adminRole := adminData[constant.JWT_ROLE].(string)
	if adminRole != constant.RoleAdmin {
		return helper.UnauthorizedError(c)
	}

	challengeID := c.Param("id")

	err = h.challengeService.Delete(challengeID)
	if err != nil {
		return c.JSON(helper.ConvertResponseCode(err), helper.FormatResponse(false, err.Error(), nil))
	}

	return c.JSON(http.StatusOK, helper.FormatResponse(true, "Challenge deleted successfully", nil))
}

// Challenge Task
// @Summary      Create a new task for a challenge
// @Description  Create a new task associated with a specific challenge. Requires admin role.
// @Tags         Challenge Tasks (Admin)
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string               true   "Bearer Token"
// @Param        body           body      ChallengeTaskRequest true   "Task data"
// @Success      201  {object}  helper.Response{data=string}   "Task created successfully"
// @Failure      400  {object}  helper.Response{data=string}   "Bad request"
// @Failure      401  {object}  helper.Response{data=string}   "Unauthorized"
// @Failure      500  {object}  helper.Response{data=string}   "Internal server error"
// @Router       /admin/challenges/tasks [post]
func (h *ChallengeHandler) CreateTask(c echo.Context) error {
	tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
	if tokenString == "" {
		return helper.UnauthorizedError(c)
	}

	token, err := h.jwt.ValidateToken(tokenString)
	if err != nil {
		return helper.UnauthorizedError(c)
	}

	adminData := h.jwt.ExtractUserToken(token)
	adminRole := adminData[constant.JWT_ROLE].(string)
	if adminRole != constant.RoleAdmin {
		return helper.UnauthorizedError(c)
	}

	var taskRequest ChallengeTaskRequest
	if err := c.Bind(&taskRequest); err != nil {
		return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "Error bad request", nil))
	}

	if err := c.Validate(taskRequest); err != nil {
		return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, err.Error(), nil))
	}

	err = h.challengeService.CreateTask(taskRequest.ChallengeID, taskRequest.DayNumber, taskRequest.TaskDescription)
	if err != nil {
		return c.JSON(helper.ConvertResponseCode(err), helper.FormatResponse(false, err.Error(), nil))
	}

	return c.JSON(http.StatusCreated, helper.FormatResponse(true, "Success", nil))
}

// @Summary      Get all tasks by challenge ID
// @Description  Retrieve all tasks associated with a specific challenge.
// @Tags         Challenge Tasks (Admin)
// @Accept       json
// @Produce      json
// @Param        challenge_id   path      string  true   "Challenge ID"
// @Success      200  {object}  helper.Response{data=[]string} "Tasks retrieved successfully"
// @Failure      404  {object}  helper.Response{data=string} "Challenge not found"
// @Failure      500  {object}  helper.Response{data=string} "Internal server error"
// @Router       /admin/challenges/{challenge_id}/tasks [get]
func (h *ChallengeHandler) GetAllTasksByChallengeID(c echo.Context) error {
	tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
	if tokenString == "" {
		return helper.UnauthorizedError(c)
	}

	token, err := h.jwt.ValidateToken(tokenString)
	if err != nil {
		return helper.UnauthorizedError(c)
	}

	adminData := h.jwt.ExtractUserToken(token)
	adminRole := adminData[constant.JWT_ROLE].(string)
	if adminRole != constant.RoleAdmin {
		return helper.UnauthorizedError(c)
	}

	challengeID := c.Param("challenge_id")
	tasks, err := h.challengeService.GetAllTasksByChallengeID(challengeID)
	if err != nil {
		return c.JSON(helper.ConvertResponseCode(err), helper.FormatResponse(false, err.Error(), nil))
	}
	return c.JSON(http.StatusOK, helper.FormatResponse(true, "Success", tasks))
}

// @Summary      Get task by ID
// @Description  Retrieve a specific task by its ID.
// @Tags         Challenge Tasks (Admin)
// @Accept       json
// @Produce      json
// @Param        task_id        path      string  true   "Task ID"
// @Success      200  {object}  helper.Response{data=string} "Task retrieved successfully"
// @Failure      404  {object}  helper.Response{data=string} "Task not found"
// @Failure      500  {object}  helper.Response{data=string} "Internal server error"
// @Router       /admin/challenges/tasks/{task_id} [get]
func (h *ChallengeHandler) GetTaskByID(c echo.Context) error {
	tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
	if tokenString == "" {
		return helper.UnauthorizedError(c)
	}

	token, err := h.jwt.ValidateToken(tokenString)
	if err != nil {
		return helper.UnauthorizedError(c)
	}

	adminData := h.jwt.ExtractUserToken(token)
	adminRole := adminData[constant.JWT_ROLE].(string)
	if adminRole != constant.RoleAdmin {
		return helper.UnauthorizedError(c)
	}

	taskID := c.Param("task_id")
	task, err := h.challengeService.GetTaskByID(taskID)
	if err != nil {
		return c.JSON(helper.ConvertResponseCode(err), helper.FormatResponse(false, err.Error(), nil))
	}
	return c.JSON(http.StatusOK, helper.FormatResponse(true, "Success", task))
}

// @Summary      Update a task
// @Description  Update the description of a specific task. Requires admin role.
// @Tags         Challenge Tasks
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string               true   "Bearer Token"
// @Param        task_id        path      string               true   "Task ID"
// @Param        body           body      ChallengeTaskRequest true   "Updated task data"
// @Success      200  {object}  helper.Response{data=string}   "Task updated successfully"
// @Failure      400  {object}  helper.Response{data=string}   "Bad request"
// @Failure      401  {object}  helper.Response{data=string}   "Unauthorized"
// @Failure      404  {object}  helper.Response{data=string}   "Task not found"
// @Failure      500  {object}  helper.Response{data=string}   "Internal server error"
// @Router       /admin/challenges/tasks/{task_id} [put]
func (h *ChallengeHandler) UpdateTask(c echo.Context) error {
	tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
	if tokenString == "" {
		return helper.UnauthorizedError(c)
	}

	token, err := h.jwt.ValidateToken(tokenString)
	if err != nil {
		return helper.UnauthorizedError(c)
	}

	adminData := h.jwt.ExtractUserToken(token)
	adminRole := adminData[constant.JWT_ROLE].(string)
	if adminRole != constant.RoleAdmin {
		return helper.UnauthorizedError(c)
	}

	var taskRequest ChallengeTaskRequest
	if err := c.Bind(&taskRequest); err != nil {
		return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "Error bad request", nil))
	}

	if err := c.Validate(taskRequest); err != nil {
		return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, err.Error(), nil))
	}

	taskID := c.Param("task_id")
	err = h.challengeService.UpdateTask(taskID, taskRequest.TaskDescription)
	if err != nil {
		return c.JSON(helper.ConvertResponseCode(err), helper.FormatResponse(false, err.Error(), nil))
	}

	return c.JSON(http.StatusOK, helper.FormatResponse(true, "Success", nil))
}

// @Summary      Delete a task
// @Description  Remove a specific task by its ID. Requires admin role.
// @Tags         Challenge Tasks
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true   "Bearer Token"
// @Param        task_id        path      string  true   "Task ID"
// @Success      200  {object}  helper.Response{data=string} "Task deleted successfully"
// @Failure      401  {object}  helper.Response{data=string} "Unauthorized"
// @Failure      404  {object}  helper.Response{data=string} "Task not found"
// @Failure      500  {object}  helper.Response{data=string} "Internal server error"
// @Router       /admin/challenges/tasks/{task_id} [delete]
func (h *ChallengeHandler) DeleteTask(c echo.Context) error {
	tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
	if tokenString == "" {
		return helper.UnauthorizedError(c)
	}

	token, err := h.jwt.ValidateToken(tokenString)
	if err != nil {
		return helper.UnauthorizedError(c)
	}

	adminData := h.jwt.ExtractUserToken(token)
	adminRole := adminData[constant.JWT_ROLE].(string)
	if adminRole != constant.RoleAdmin {
		return helper.UnauthorizedError(c)
	}

	taskID := c.Param("task_id")
	err = h.challengeService.DeleteTask(taskID)
	if err != nil {
		return c.JSON(helper.ConvertResponseCode(err), helper.FormatResponse(false, err.Error(), nil))
	}

	return c.JSON(http.StatusOK, helper.FormatResponse(true, "Task deleted successfully", nil))
}

// User
// CreateChallengeLog
// @Summary      Log a user's participation in a challenge
// @Description  Logs challenge participation with status "Progress"
// @Tags         Challenge Logs
// @Accept       json
// @Produce      json
// @Param        Authorization header    string               true  "Bearer Token"
// @Param        request       body      ChallengeLogRequest  true  "Challenge Log Request"
// @Success      201          {object}  helper.Response{data=string}  "Challenge log created successfully"
// @Failure      400          {object}  helper.Response{data=string}  "Bad request"
// @Failure      401          {object}  helper.Response{data=string}  "Unauthorized"
// @Failure      500          {object}  helper.Response{data=string}  "Internal server error"
// @Router       /challenges/logs [post]
func (h *ChallengeHandler) CreateChallengeLog(c echo.Context) error {
	tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
	if tokenString == "" {
		return helper.UnauthorizedError(c)
	}

	token, err := h.jwt.ValidateToken(tokenString)
	if err != nil {
		return helper.UnauthorizedError(c)
	}

	userData := h.jwt.ExtractUserToken(token)
	userID := userData[constant.JWT_ID].(string)

	var req ChallengeLogRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "Invalid request format", nil))
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, err.Error(), nil))
	}

	log := challenges.ChallengeLog{
		ChallengeID: req.ChallengeID,
		UserID:      userID,
		Status:      "Progress",
		StartDate:   time.Now(),
		Feed:        req.Feed,
	}

	err = h.challengeService.CreateChallengeLogWithConfirmation(log)
	if err != nil {
		return c.JSON(helper.ConvertResponseCode(err), helper.FormatResponse(false, err.Error(), nil))
	}

	return c.JSON(http.StatusCreated, helper.FormatResponse(true, "Challenge log created successfully", nil))
}

// UpdateChallengeConfirmationProgress
// @Summary      Update the progress of a challenge confirmation
// @Description  Updates the status of a challenge confirmation to "Done" and uploads the confirmation image
// @Tags         Challenge Confirmations
// @Accept       multipart/form-data
// @Produce      json
// @Param        Authorization         header    string                            true  "Bearer Token"
// @Param        request               formData  ChallengeConfirmationRequest      true  "Challenge Confirmation Request"
// @Param        challenge_confirmation_img formData file                          true  "Challenge Confirmation Image"
// @Success      200                   {object}  helper.Response{data=string}      "Confirmation updated successfully"
// @Failure      400                   {object}  helper.Response{data=string}      "Bad request"
// @Failure      401                   {object}  helper.Response{data=string}      "Unauthorized"
// @Failure      500                   {object}  helper.Response{data=string}      "Internal server error"
// @Router       /challenges/confirmations/progress [put]
func (h *ChallengeHandler) UpdateChallengeConfirmationProgress(c echo.Context) error {
	tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
	if tokenString == "" {
		return helper.UnauthorizedError(c)
	}

	token, err := h.jwt.ValidateToken(tokenString)
	if err != nil {
		return helper.UnauthorizedError(c)
	}

	userData := h.jwt.ExtractUserToken(token)
	userID := userData[constant.JWT_ID].(string)

	var req ChallengeConfirmationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "Invalid request format", nil))
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, err.Error(), nil))
	}

	file, err := c.FormFile("challenge_confirmation_img")
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "Challenge confirmation image is required", nil))
	}

	src, err := h.storage.ImageValidation(file)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, err.Error(), nil))
	}

	challengeImgURL, err := h.storage.UploadImageToCloudinary(src, "ecomate/challenges/confirmations/images/")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, "Failed to upload challenge confirmation image", nil))
	}

	err = h.challengeService.UpdateChallengeConfirmationProgress(req.ChallengeConfirmationID, challengeImgURL, userID)
	if err != nil {
		return c.JSON(helper.ConvertResponseCode(err), helper.FormatResponse(false, err.Error(), nil))
	}

	err = h.challengeService.CheckAndUpdateChallengeLogStatusByConfirmation(req.ChallengeConfirmationID, userID)
	if err != nil {
		return c.JSON(helper.ConvertResponseCode(err), helper.FormatResponse(false, err.Error(), nil))
	}

	return c.JSON(http.StatusOK, helper.FormatResponse(true, "Confirmation updated successfully", nil))
}

// ClaimRewards
// @Summary      Claim rewards for a challenge
// @Description  Claims exp and coin rewards from a completed challenge
// @Tags         Challenge Rewards
// @Accept       json
// @Produce      json
// @Param        Authorization header    string                     true  "Bearer Token"
// @Param        request       body      ClaimRewardsRequest        true  "Claim Rewards Request"
// @Success      200          {object}  helper.Response{data=string}  "Rewards claimed successfully"
// @Failure      400          {object}  helper.Response{data=string}  "Bad request"
// @Failure      401          {object}  helper.Response{data=string}  "Unauthorized"
// @Failure      404          {object}  helper.Response{data=string}  "Challenge log not found"
// @Failure      500          {object}  helper.Response{data=string}  "Internal server error"
// @Router       /challenges/rewards [post]
func (h *ChallengeHandler) ClaimRewards(c echo.Context) error {
	tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
	if tokenString == "" {
		return helper.UnauthorizedError(c)
	}

	token, err := h.jwt.ValidateToken(tokenString)
	if err != nil {
		return helper.UnauthorizedError(c)
	}

	userData := h.jwt.ExtractUserToken(token)
	userID := userData[constant.JWT_ID].(string)

	var req ClaimRewardsRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "Invalid request format", nil))
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, err.Error(), nil))
	}

	err = h.challengeService.ClaimRewards(req.ChallengeLogID, userID)
	if err != nil {
		if errors.Is(err, constant.ErrRewardAlreadyClaimed) {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "Reward already claimed", nil))
		}
		return c.JSON(helper.ConvertResponseCode(err), helper.FormatResponse(false, err.Error(), nil))
	}

	return c.JSON(http.StatusOK, helper.FormatResponse(true, "Rewards claimed successfully", nil))
}

// GetActiveChallenges retrieves active challenges for the user
// @Summary      Get active challenges
// @Description  Get challenges that are currently in progress for the user
// @Tags         Challenges (User)
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer Token"
// @Success      200 {object} helper.Response{data=[]challenges.ChallengeLog} "Active challenges retrieved successfully"
// @Failure      401 {object} helper.Response{data=string} "Unauthorized"
// @Failure      500 {object} helper.Response{data=string} "Internal server error"
// @Router       /challenges/active [get]
func (h *ChallengeHandler) GetActiveChallenges(c echo.Context) error {
	tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
	if tokenString == "" {
		return helper.UnauthorizedError(c)
	}

	token, err := h.jwt.ValidateToken(tokenString)
	if err != nil {
		return helper.UnauthorizedError(c)
	}

	userData := h.jwt.ExtractUserToken(token)
	userID := userData[constant.JWT_ID].(string)

	challenges, err := h.challengeService.GetActiveChallenges(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, err.Error(), nil))
	}

	if len(challenges) == 0 {
		return c.NoContent(http.StatusNoContent)
	}

	return c.JSON(http.StatusOK, helper.FormatResponse(true, "Active challenges retrieved successfully", challenges))
}

// GetUnclaimedChallenges retrieves challenges not yet taken by the user
// @Summary      Get unclaimed challenges
// @Description  Get challenges that the user has not taken yet
// @Tags         Challenges (User)
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer Token"
// @Success      200 {object} helper.Response{data=[]challenges.Challenge} "Unclaimed challenges retrieved successfully"
// @Failure      401 {object} helper.Response{data=string} "Unauthorized"
// @Failure      500 {object} helper.Response{data=string} "Internal server error"
// @Router       /challenges/unclaimed [get]
func (h *ChallengeHandler) GetUnclaimedChallenges(c echo.Context) error {
	tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
	if tokenString == "" {
		return helper.UnauthorizedError(c)
	}

	token, err := h.jwt.ValidateToken(tokenString)
	if err != nil {
		return helper.UnauthorizedError(c)
	}

	userData := h.jwt.ExtractUserToken(token)
	userID := userData[constant.JWT_ID].(string)
	role := userData[constant.JWT_ROLE].(string)

	isAdmin := role == "admin"

	challenges, err := h.challengeService.GetUnclaimedChallenges(userID, isAdmin)
	if err != nil {
		return c.JSON(helper.ConvertResponseCode(err), helper.FormatResponse(false, err.Error(), nil))
	}

	if len(challenges) == 0 {
		return c.JSON(http.StatusNotFound, helper.FormatResponse(false, "No unclaimed challenges available", nil))
	}

	return c.JSON(http.StatusOK, helper.FormatResponse(true, "Unclaimed challenges retrieved successfully", challenges))
}

// GetChallengeDetailsWithConfirmations retrieves challenge log and its confirmations
// @Summary      Get challenge details with confirmations
// @Description  Get challenge log and confirmation tasks for a specific challenge log ID
// @Tags         Challenges (User)
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer Token"
// @Param        challengeLogID query string true "Challenge Log ID"
// @Success      200 {object} helper.Response{data=challenges.ChallengeLogDetails} "Challenge details retrieved successfully"
// @Failure      401 {object} helper.Response{data=string} "Unauthorized"
// @Failure      404 {object} helper.Response{data=string} "Challenge log not found"
// @Failure      500 {object} helper.Response{data=string} "Internal server error"
// @Router       /challenges/details [get]
func (h *ChallengeHandler) GetChallengeDetailsWithConfirmations(c echo.Context) error {
	tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
	if tokenString == "" {
		return helper.UnauthorizedError(c)
	}

	token, err := h.jwt.ValidateToken(tokenString)
	if err != nil {
		return helper.UnauthorizedError(c)
	}

	userData := h.jwt.ExtractUserToken(token)
	userID := userData[constant.JWT_ID].(string)

	challengeLogID := c.QueryParam("challengeLogID")
	if challengeLogID == "" {
		return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "Challenge log ID is required", nil))
	}

	details, err := h.challengeService.GetChallengeDetailsWithConfirmations(userID, challengeLogID)
	if err != nil {
		if errors.Is(err, constant.ErrUnauthorized) {
			return helper.UnauthorizedError(c)
		}
		if errors.Is(err, constant.ErrChallengeLogNotFound) {
			return c.JSON(http.StatusNotFound, helper.FormatResponse(false, err.Error(), nil))
		}
		return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, err.Error(), nil))
	}

	return c.JSON(http.StatusOK, helper.FormatResponse(true, "Challenge details retrieved successfully", details))
}

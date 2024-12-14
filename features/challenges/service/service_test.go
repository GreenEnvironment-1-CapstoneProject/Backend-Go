package service

import (
	"errors"
	"greenenvironment/constant"
	"greenenvironment/features/challenges"
	"greenenvironment/features/impacts"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockChallengeRepository struct {
	mock.Mock
}

func (m *MockChallengeRepository) Create(challenge challenges.Challenge) error {
	args := m.Called(challenge)
	return args.Error(0)
}

func (m *MockChallengeRepository) GetAllByPage(page int) ([]challenges.Challenge, int, error) {
	args := m.Called(page)
	return args.Get(0).([]challenges.Challenge), args.Int(1), args.Error(2)
}

func (m *MockChallengeRepository) GetByID(id string) (challenges.Challenge, error) {
	args := m.Called(id)
	return args.Get(0).(challenges.Challenge), args.Error(1)
}

func (m *MockChallengeRepository) Update(challengeData challenges.Challenge) error {
	args := m.Called(challengeData)
	return args.Error(0)
}

func (m *MockChallengeRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockChallengeRepository) CreateTask(task challenges.ChallengeTask) error {
	args := m.Called(task)
	return args.Error(0)
}

func (m *MockChallengeRepository) GetTasksByChallengeID(challengeID string) ([]challenges.ChallengeTask, error) {
	args := m.Called(challengeID)
	return args.Get(0).([]challenges.ChallengeTask), args.Error(1)
}

func (m *MockChallengeRepository) GetTaskByID(taskID string) (challenges.ChallengeTask, error) {
	args := m.Called(taskID)
	return args.Get(0).(challenges.ChallengeTask), args.Error(1)
}

func (m *MockChallengeRepository) UpdateTask(task challenges.ChallengeTask) error {
	args := m.Called(task)
	return args.Error(0)
}

func (m *MockChallengeRepository) DeleteTask(taskID string) error {
	args := m.Called(taskID)
	return args.Error(0)
}

func (m *MockChallengeRepository) CreateChallengeLog(log challenges.ChallengeLog) error {
	args := m.Called(log)
	return args.Error(0)
}

func (m *MockChallengeRepository) CreateChallengeConfirmation(confirmation challenges.ChallengeConfirmation) error {
	args := m.Called(confirmation)
	return args.Error(0)
}

func (m *MockChallengeRepository) IncrementChallengeCounts(challengeID string, actionCount int, participantIncrement bool) error {
	args := m.Called(challengeID, actionCount, participantIncrement)
	return args.Error(0)
}

func (m *MockChallengeRepository) IsChallengeTaken(userID, challengeID string) (bool, error) {
	args := m.Called(userID, challengeID)
	return args.Bool(0), args.Error(1)
}

func (m *MockChallengeRepository) GetChallengeConfirmationByID(confirmationID string) (challenges.ChallengeConfirmation, error) {
	args := m.Called(confirmationID)
	return args.Get(0).(challenges.ChallengeConfirmation), args.Error(1)
}

func (m *MockChallengeRepository) UpdateChallengeConfirmation(confirmation challenges.ChallengeConfirmation) error {
	args := m.Called(confirmation)
	return args.Error(0)
}

func (m *MockChallengeRepository) GetChallengeTaskByID(taskID string) (challenges.ChallengeTask, error) {
	args := m.Called(taskID)
	return args.Get(0).(challenges.ChallengeTask), args.Error(1)
}

func (m *MockChallengeRepository) GetChallengeLogByChallengeIDAndUserID(challengeID, userID string) (challenges.ChallengeLog, error) {
	args := m.Called(challengeID, userID)
	return args.Get(0).(challenges.ChallengeLog), args.Error(1)
}

func (m *MockChallengeRepository) GetConfirmationsByChallengeID(challengeID, userID string) ([]challenges.ChallengeConfirmation, error) {
	args := m.Called(challengeID, userID)
	return args.Get(0).([]challenges.ChallengeConfirmation), args.Error(1)
}

func (m *MockChallengeRepository) UpdateChallengeLog(log challenges.ChallengeLog) error {
	args := m.Called(log)
	return args.Error(0)
}

func (m *MockChallengeRepository) UpdateTaskAndChallengeStatus() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockChallengeRepository) IsRewardClaimed(challengeLogID string) (bool, error) {
	args := m.Called(challengeLogID)
	return args.Bool(0), args.Error(1)
}

func (m *MockChallengeRepository) UpdateRewardsGiven(challengeLogID string) error {
	args := m.Called(challengeLogID)
	return args.Error(0)
}

func (m *MockChallengeRepository) AddUserRewards(userID string, exp int, coin int) error {
	args := m.Called(userID, exp, coin)
	return args.Error(0)
}

func (m *MockChallengeRepository) GetChallengeIDByLogID(challengeLogID string) (string, error) {
	args := m.Called(challengeLogID)
	return args.String(0), args.Error(1)
}

func (m *MockChallengeRepository) GetChallengeRewards(challengeID string) (int, int, error) {
	args := m.Called(challengeID)
	return args.Int(0), args.Int(1), args.Error(2)
}

func (m *MockChallengeRepository) GetChallengeLogByUserID(userID string, page, perPage int, difficulty, title string) ([]challenges.ChallengeLog, int, error) {
	args := m.Called(userID, page, perPage, difficulty, title)
	return args.Get(0).([]challenges.ChallengeLog), args.Int(1), args.Error(2)
}

func (m *MockChallengeRepository) GetUnclaimedChallenges(userID string, isAdmin bool, page, limit int, difficulty, title string) ([]challenges.Challenge, int, error) {
	args := m.Called(userID, isAdmin, page, limit, difficulty, title)
	return args.Get(0).([]challenges.Challenge), args.Int(1), args.Error(2)
}

func (m *MockChallengeRepository) GetChallengeLogByID(challengeLogID string) (challenges.ChallengeLog, error) {
	args := m.Called(challengeLogID)
	return args.Get(0).(challenges.ChallengeLog), args.Error(1)
}

func (m *MockChallengeRepository) GetChallengeByID(challengeID string) (challenges.Challenge, error) {
	args := m.Called(challengeID)
	return args.Get(0).(challenges.Challenge), args.Error(1)
}

func (m *MockChallengeRepository) GetTasksByChallengeIDforUser(challengeID string) ([]challenges.ChallengeTask, error) {
	args := m.Called(challengeID)
	return args.Get(0).([]challenges.ChallengeTask), args.Error(1)
}

type MockImpactRepository struct {
	mock.Mock
}

func (m *MockImpactRepository) GetAll() ([]impacts.ImpactCategory, error) {
	args := m.Called()
	return args.Get(0).([]impacts.ImpactCategory), args.Error(1)
}

func (m *MockImpactRepository) GetByID(ID string) (impacts.ImpactCategory, error) {
	args := m.Called(ID)
	return args.Get(0).(impacts.ImpactCategory), args.Error(1)
}

func (m *MockImpactRepository) Create(category impacts.ImpactCategory) error {
	args := m.Called(category)
	return args.Error(0)
}

func (m *MockImpactRepository) Delete(category impacts.ImpactCategory) error {
	args := m.Called(category)
	return args.Error(0)
}

func TestCreateChallenge_Success(t *testing.T) {
	mockChallengeRepo := new(MockChallengeRepository)
	mockImpactRepo := new(MockImpactRepository)
	service := NewChallengeService(mockChallengeRepo, mockImpactRepo)

	newChallenge := challenges.Challenge{
		Title:       "New Challenge",
		Difficulty:  "Easy",
		Description: "Test Description",
		ImpactCategories: []challenges.ChallengeImpactCategory{
			{ImpactCategoryID: "impact1"},
		},
	}

	mockImpactRepo.On("GetByID", "impact1").Return(impacts.ImpactCategory{
		ID:          "impact1",
		Name:        "Environment",
		ImpactPoint: 10,
		Description: "Environmental Impact",
	}, nil)

	mockChallengeRepo.On("Create", mock.Anything).Return(nil)

	err := service.Create(newChallenge)

	assert.NoError(t, err)
	mockImpactRepo.AssertCalled(t, "GetByID", "impact1")
	mockChallengeRepo.AssertCalled(t, "Create", mock.Anything)
}

func TestCreateChallenge_InvalidImpactCategory(t *testing.T) {
	mockChallengeRepo := new(MockChallengeRepository)
	mockImpactRepo := new(MockImpactRepository)
	service := NewChallengeService(mockChallengeRepo, mockImpactRepo)

	newChallenge := challenges.Challenge{
		Title: "New Challenge",
		ImpactCategories: []challenges.ChallengeImpactCategory{
			{ImpactCategoryID: "invalidImpact"},
		},
	}

	mockImpactRepo.On("GetByID", "invalidImpact").Return(impacts.ImpactCategory{}, nil)

	err := service.Create(newChallenge)

	assert.Error(t, err)
	assert.Equal(t, constant.ErrCreateProduct, err)
	mockImpactRepo.AssertCalled(t, "GetByID", "invalidImpact")
	mockChallengeRepo.AssertNotCalled(t, "Create", mock.Anything)
}

func TestCreateChallenge_RepositoryError(t *testing.T) {
	mockChallengeRepo := new(MockChallengeRepository)
	mockImpactRepo := new(MockImpactRepository)
	service := NewChallengeService(mockChallengeRepo, mockImpactRepo)

	newChallenge := challenges.Challenge{
		Title:       "New Challenge",
		Difficulty:  "Easy",
		Description: "Test Description",
		ImpactCategories: []challenges.ChallengeImpactCategory{
			{ImpactCategoryID: "impact1"},
		},
	}

	mockImpactRepo.On("GetByID", "impact1").Return(impacts.ImpactCategory{
		ID: "impact1",
	}, nil)

	mockChallengeRepo.On("Create", mock.Anything).Return(errors.New("database error"))

	err := service.Create(newChallenge)

	assert.Error(t, err)
	assert.Equal(t, "database error", err.Error())
	mockImpactRepo.AssertCalled(t, "GetByID", "impact1")
	mockChallengeRepo.AssertCalled(t, "Create", mock.Anything)
}

func TestGetChallengeByID_Success(t *testing.T) {
	mockChallengeRepo := new(MockChallengeRepository)
	service := NewChallengeService(mockChallengeRepo, nil)

	expectedChallenge := challenges.Challenge{
		ID:          "challenge1",
		Title:       "Challenge Title",
		Description: "Challenge Description",
	}

	mockChallengeRepo.On("GetByID", "challenge1").Return(expectedChallenge, nil)

	challenge, err := service.GetByID("challenge1")

	assert.NoError(t, err)
	assert.Equal(t, "challenge1", challenge.ID)
	assert.Equal(t, "Challenge Title", challenge.Title)
	mockChallengeRepo.AssertCalled(t, "GetByID", "challenge1")
}

func TestGetChallengeByID_NotFound(t *testing.T) {
	mockChallengeRepo := new(MockChallengeRepository)
	service := NewChallengeService(mockChallengeRepo, nil)

	mockChallengeRepo.On("GetByID", "invalidID").Return(challenges.Challenge{}, errors.New("not found"))

	challenge, err := service.GetByID("invalidID")

	assert.Error(t, err)
	assert.Equal(t, "not found", err.Error())
	assert.Equal(t, "", challenge.ID)
	mockChallengeRepo.AssertCalled(t, "GetByID", "invalidID")
}

func TestGetAllByPage_Success(t *testing.T) {
	mockChallengeRepo := new(MockChallengeRepository)
	service := NewChallengeService(mockChallengeRepo, nil)

	mockChallengeRepo.On("GetAllByPage", 1).Return([]challenges.Challenge{
		{
			ID:          "challenge1",
			Title:       "Challenge 1",
			Difficulty:  "Medium",
			Description: "Test Description",
		},
	}, 2, nil)

	challenges, totalPages, err := service.GetAllByPage(1)

	assert.NoError(t, err)
	assert.Equal(t, 2, totalPages)
	assert.Len(t, challenges, 1)
	assert.Equal(t, "Challenge 1", challenges[0].Title)
	mockChallengeRepo.AssertCalled(t, "GetAllByPage", 1)
}

func TestUpdateChallenge_Success(t *testing.T) {
	mockChallengeRepo := new(MockChallengeRepository)
	mockImpactRepo := new(MockImpactRepository)
	service := NewChallengeService(mockChallengeRepo, mockImpactRepo)

	updatedChallenge := challenges.Challenge{
		ID:          "challenge1",
		Title:       "Updated Title",
		Description: "Updated Description",
		ImpactCategories: []challenges.ChallengeImpactCategory{
			{ImpactCategoryID: "impact1"},
		},
	}

	existingChallenge := challenges.Challenge{
		ID:           "challenge1",
		Author:       "original-author",
		ChallengeImg: "original-image.png",
	}

	mockImpactRepo.On("GetByID", "impact1").Return(impacts.ImpactCategory{
		ID:          "impact1",
		Name:        "Environment",
		ImpactPoint: 10,
		Description: "Environmental Impact",
	}, nil)

	mockChallengeRepo.On("GetByID", "challenge1").Return(existingChallenge, nil)

	mockChallengeRepo.On("Update", mock.MatchedBy(func(challenge challenges.Challenge) bool {
		assert.Equal(t, "original-author", challenge.Author)
		assert.Equal(t, "original-image.png", challenge.ChallengeImg)
		return true
	})).Return(nil)

	err := service.Update(updatedChallenge)

	assert.NoError(t, err)
	mockImpactRepo.AssertCalled(t, "GetByID", "impact1")
	mockChallengeRepo.AssertCalled(t, "GetByID", "challenge1")
	mockChallengeRepo.AssertCalled(t, "Update", mock.Anything)
}

func TestUpdateChallenge_InvalidImpactCategory(t *testing.T) {
	mockChallengeRepo := new(MockChallengeRepository)
	mockImpactRepo := new(MockImpactRepository)
	service := NewChallengeService(mockChallengeRepo, mockImpactRepo)

	invalidChallenge := challenges.Challenge{
		ID:    "challenge1",
		Title: "Invalid Challenge",
		ImpactCategories: []challenges.ChallengeImpactCategory{
			{ImpactCategoryID: "invalidImpact"},
		},
	}

	mockImpactRepo.On("GetByID", "invalidImpact").Return(impacts.ImpactCategory{}, nil)

	err := service.Update(invalidChallenge)

	assert.Error(t, err)
	assert.Equal(t, constant.ErrUpdateChallenge, err)
	mockImpactRepo.AssertCalled(t, "GetByID", "invalidImpact")
	mockChallengeRepo.AssertNotCalled(t, "Update", mock.Anything)
}

func TestUpdateChallenge_RepositoryError(t *testing.T) {
	mockChallengeRepo := new(MockChallengeRepository)
	mockImpactRepo := new(MockImpactRepository)
	service := NewChallengeService(mockChallengeRepo, mockImpactRepo)

	updateChallenge := challenges.Challenge{
		ID:    "challenge1",
		Title: "Update Error Challenge",
		ImpactCategories: []challenges.ChallengeImpactCategory{
			{ImpactCategoryID: "impact1"},
		},
	}

	existingChallenge := challenges.Challenge{
		ID: "challenge1",
	}

	mockChallengeRepo.On("GetByID", "challenge1").Return(existingChallenge, nil)
	mockImpactRepo.On("GetByID", "impact1").Return(impacts.ImpactCategory{ID: "impact1"}, nil)
	mockChallengeRepo.On("Update", mock.Anything).Return(errors.New("repository error"))

	err := service.Update(updateChallenge)

	assert.Error(t, err)
	assert.Equal(t, "repository error", err.Error())
	mockChallengeRepo.AssertCalled(t, "Update", mock.Anything)
}

func TestDeleteChallenge_Success(t *testing.T) {
	mockChallengeRepo := new(MockChallengeRepository)
	service := NewChallengeService(mockChallengeRepo, nil)

	mockChallengeRepo.On("GetByID", "challenge1").Return(challenges.Challenge{ID: "challenge1"}, nil)
	mockChallengeRepo.On("Delete", "challenge1").Return(nil)

	err := service.Delete("challenge1")

	assert.NoError(t, err)
	mockChallengeRepo.AssertCalled(t, "GetByID", "challenge1")
	mockChallengeRepo.AssertCalled(t, "Delete", "challenge1")
}

func TestDeleteChallenge_NotFound(t *testing.T) {
	mockChallengeRepo := new(MockChallengeRepository)
	service := NewChallengeService(mockChallengeRepo, nil)

	mockChallengeRepo.On("GetByID", "challenge1").Return(challenges.Challenge{}, constant.ErrChallengeNotFound)

	err := service.Delete("challenge1")

	assert.Error(t, err)
	assert.Equal(t, constant.ErrChallengeNotFound, err)
	mockChallengeRepo.AssertCalled(t, "GetByID", "challenge1")
	mockChallengeRepo.AssertNotCalled(t, "Delete", mock.Anything)
}

func TestCreateTask_Success(t *testing.T) {
	mockChallengeRepo := new(MockChallengeRepository)
	mockImpactRepo := new(MockImpactRepository)
	service := NewChallengeService(mockChallengeRepo, mockImpactRepo)

	challengeID := "challenge1"
	taskName := "Task Name"
	dayNumber := 1
	taskDescription := "Task Description"

	mockChallengeRepo.On("GetByID", challengeID).Return(challenges.Challenge{
		ID:           challengeID,
		DurationDays: 5,
	}, nil)

	mockChallengeRepo.On("GetTasksByChallengeID", challengeID).Return([]challenges.ChallengeTask{}, nil)

	mockChallengeRepo.On("CreateTask", mock.Anything).Return(nil)

	err := service.CreateTask(challengeID, taskName, dayNumber, taskDescription)

	assert.NoError(t, err)
	mockChallengeRepo.AssertCalled(t, "GetByID", challengeID)
	mockChallengeRepo.AssertCalled(t, "GetTasksByChallengeID", challengeID)
	mockChallengeRepo.AssertCalled(t, "CreateTask", mock.Anything)
}

func TestCreateTask_ErrorOnGetByID(t *testing.T) {
	mockChallengeRepo := new(MockChallengeRepository)
	mockImpactRepo := new(MockImpactRepository)
	service := NewChallengeService(mockChallengeRepo, mockImpactRepo)

	challengeID := "challenge1"
	taskName := "Task Name"
	dayNumber := 1
	taskDescription := "Task Description"

	mockChallengeRepo.On("GetByID", challengeID).Return(challenges.Challenge{}, constant.ErrChallengeNotFound)

	err := service.CreateTask(challengeID, taskName, dayNumber, taskDescription)

	assert.Error(t, err)
	assert.EqualError(t, err, constant.ErrChallengeNotFound.Error())
	mockChallengeRepo.AssertCalled(t, "GetByID", challengeID)
}

func TestCreateTask_ChallengeNotFound(t *testing.T) {
	mockChallengeRepo := new(MockChallengeRepository)
	mockImpactRepo := new(MockImpactRepository)
	service := NewChallengeService(mockChallengeRepo, mockImpactRepo)

	challengeID := "challenge1"
	taskName := "Task Name"
	dayNumber := 1
	taskDescription := "Task Description"

	mockChallengeRepo.On("GetByID", challengeID).Return(challenges.Challenge{}, nil)

	err := service.CreateTask(challengeID, taskName, dayNumber, taskDescription)

	assert.Error(t, err)
	assert.Equal(t, constant.ErrChallengeNotFound, err)
	mockChallengeRepo.AssertCalled(t, "GetByID", challengeID)
}

func TestCreateTask_TaskAlreadyExists(t *testing.T) {
	mockChallengeRepo := new(MockChallengeRepository)
	mockImpactRepo := new(MockImpactRepository)
	service := NewChallengeService(mockChallengeRepo, mockImpactRepo)

	challengeID := "challenge1"
	existingTask := challenges.ChallengeTask{
		ID:              "task1",
		ChallengeID:     challengeID,
		Name:            "Existing Task",
		DayNumber:       1,
		TaskDescription: "Existing Task Description",
	}

	mockChallengeRepo.On("GetByID", challengeID).Return(challenges.Challenge{
		ID:           challengeID,
		DurationDays: 5,
	}, nil)

	mockChallengeRepo.On("GetTasksByChallengeID", challengeID).Return([]challenges.ChallengeTask{existingTask}, nil)

	err := service.CreateTask(challengeID, "New Task", 1, "New Task Description")

	assert.Error(t, err)
	assert.Equal(t, constant.ErrTaskAlreadyExists, err)
	mockChallengeRepo.AssertCalled(t, "GetByID", challengeID)
	mockChallengeRepo.AssertCalled(t, "GetTasksByChallengeID", challengeID)
}

func TestCreateTask_DayNumberInvalid(t *testing.T) {
	mockChallengeRepo := new(MockChallengeRepository)
	mockImpactRepo := new(MockImpactRepository)
	service := NewChallengeService(mockChallengeRepo, mockImpactRepo)

	challengeID := "challenge1"

	mockChallengeRepo.On("GetByID", challengeID).Return(challenges.Challenge{
		ID:           challengeID,
		DurationDays: 5,
	}, nil)

	err := service.CreateTask(challengeID, "Task Name", 6, "Task Description")

	assert.Error(t, err)
	assert.Equal(t, constant.ErrInvalidDayNumber, err)
	mockChallengeRepo.AssertCalled(t, "GetByID", challengeID)
}

func TestGetAllTasksByChallengeID(t *testing.T) {
	mockRepo := new(MockChallengeRepository)
	service := NewChallengeService(mockRepo, nil)

	challengeID := "challenge1"
	expectedTasks := []challenges.ChallengeTask{
		{ID: "task1", TaskDescription: "Task 1"},
		{ID: "task2", TaskDescription: "Task 2"},
	}

	mockRepo.On("GetTasksByChallengeID", challengeID).Return(expectedTasks, nil)

	tasks, err := service.GetAllTasksByChallengeID(challengeID)

	assert.NoError(t, err)
	assert.Equal(t, expectedTasks, tasks)
	mockRepo.AssertCalled(t, "GetTasksByChallengeID", challengeID)
}

func TestGetTaskByID(t *testing.T) {
	mockRepo := new(MockChallengeRepository)
	service := NewChallengeService(mockRepo, nil)

	taskID := "task1"
	existingTask := challenges.ChallengeTask{ID: taskID, TaskDescription: "Task 1"}

	mockRepo.On("GetTaskByID", taskID).Return(existingTask, nil)

	task, err := service.GetTaskByID(taskID)

	assert.NoError(t, err)
	assert.Equal(t, existingTask, task)
	mockRepo.AssertCalled(t, "GetTaskByID", taskID)
}

func TestGetTaskByID_NotFound(t *testing.T) {
	mockRepo := new(MockChallengeRepository)
	service := NewChallengeService(mockRepo, nil)

	taskID := "nonexistent"

	mockRepo.On("GetTaskByID", taskID).Return(challenges.ChallengeTask{}, errors.New("not found"))

	task, err := service.GetTaskByID(taskID)

	assert.Error(t, err)
	assert.ErrorIs(t, err, constant.ErrTaskNotFound)
	assert.Empty(t, task)
	mockRepo.AssertCalled(t, "GetTaskByID", taskID)
}

func TestUpdateTask(t *testing.T) {
	mockRepo := new(MockChallengeRepository)
	service := NewChallengeService(mockRepo, nil)

	taskID := "task1"
	taskDescription := "Updated Task Description"

	existingTask := challenges.ChallengeTask{ID: taskID, TaskDescription: "Old Description"}
	updatedTask := challenges.ChallengeTask{ID: taskID, TaskDescription: taskDescription}

	mockRepo.On("GetTaskByID", taskID).Return(existingTask, nil)
	mockRepo.On("UpdateTask", updatedTask).Return(nil)

	err := service.UpdateTask(taskID, taskDescription)

	assert.NoError(t, err)
	mockRepo.AssertCalled(t, "GetTaskByID", taskID)
	mockRepo.AssertCalled(t, "UpdateTask", mock.MatchedBy(func(task challenges.ChallengeTask) bool {
		return task.ID == updatedTask.ID && task.TaskDescription == updatedTask.TaskDescription
	}))
}

func TestUpdateTask_NotFound(t *testing.T) {
	mockRepo := new(MockChallengeRepository)
	service := NewChallengeService(mockRepo, nil)

	taskID := "nonexistent"
	taskDescription := "Updated Task Description"

	mockRepo.On("GetTaskByID", taskID).Return(challenges.ChallengeTask{}, errors.New("not found"))

	err := service.UpdateTask(taskID, taskDescription)

	assert.Error(t, err)
	assert.ErrorIs(t, err, constant.ErrTaskNotFound)
	mockRepo.AssertCalled(t, "GetTaskByID", taskID)
}

func TestDeleteTask(t *testing.T) {
	mockRepo := new(MockChallengeRepository)
	service := NewChallengeService(mockRepo, nil)

	taskID := "task1"
	existingTask := challenges.ChallengeTask{ID: taskID, TaskDescription: "Task to delete"}

	mockRepo.On("GetTaskByID", taskID).Return(existingTask, nil)
	mockRepo.On("DeleteTask", taskID).Return(nil)

	err := service.DeleteTask(taskID)

	assert.NoError(t, err)
	mockRepo.AssertCalled(t, "GetTaskByID", taskID)
	mockRepo.AssertCalled(t, "DeleteTask", taskID)
}

func TestDeleteTask_NotFound(t *testing.T) {
	mockRepo := new(MockChallengeRepository)
	service := NewChallengeService(mockRepo, nil)

	taskID := "nonexistent"

	mockRepo.On("GetTaskByID", taskID).Return(challenges.ChallengeTask{}, errors.New("not found"))

	err := service.DeleteTask(taskID)

	assert.Error(t, err)
	assert.ErrorIs(t, err, constant.ErrTaskNotFound)
	mockRepo.AssertCalled(t, "GetTaskByID", taskID)
}

func TestCreateChallengeLogWithConfirmation_Success(t *testing.T) {
	mockChallengeRepo := new(MockChallengeRepository)
	mockImpactRepo := new(MockImpactRepository)
	service := NewChallengeService(mockChallengeRepo, mockImpactRepo)

	challengeID := "challenge1"
	userID := "user1"
	log := challenges.ChallengeLog{
		ChallengeID: challengeID,
		UserID:      userID,
	}

	mockChallengeRepo.On("IsChallengeTaken", userID, challengeID).Return(false, nil)
	mockChallengeRepo.On("CreateChallengeLog", mock.Anything).Return(nil)

	mockChallengeRepo.On("GetTasksByChallengeID", challengeID).Return([]challenges.ChallengeTask{
		{ID: "task1"},
		{ID: "task2"},
	}, nil)

	mockChallengeRepo.On("CreateChallengeConfirmation", mock.Anything).Return(nil).Twice()
	mockChallengeRepo.On("IncrementChallengeCounts", challengeID, 2, true).Return(nil)

	err := service.CreateChallengeLogWithConfirmation(log)

	assert.NoError(t, err)
	mockChallengeRepo.AssertCalled(t, "IsChallengeTaken", userID, challengeID)
	mockChallengeRepo.AssertCalled(t, "CreateChallengeLog", mock.Anything)
	mockChallengeRepo.AssertCalled(t, "GetTasksByChallengeID", challengeID)
	mockChallengeRepo.AssertCalled(t, "CreateChallengeConfirmation", mock.Anything)
	mockChallengeRepo.AssertCalled(t, "IncrementChallengeCounts", challengeID, 2, true)
}

func TestCreateChallengeLogWithConfirmation_AlreadyTaken(t *testing.T) {
	mockChallengeRepo := new(MockChallengeRepository)
	mockImpactRepo := new(MockImpactRepository)
	service := NewChallengeService(mockChallengeRepo, mockImpactRepo)

	challengeID := "challenge1"
	userID := "user1"
	log := challenges.ChallengeLog{
		ChallengeID: challengeID,
		UserID:      userID,
	}

	mockChallengeRepo.On("IsChallengeTaken", userID, challengeID).Return(true, nil)

	err := service.CreateChallengeLogWithConfirmation(log)

	assert.Error(t, err)
	assert.Equal(t, constant.ErrChallengeAlreadyTaken, err)
	mockChallengeRepo.AssertCalled(t, "IsChallengeTaken", userID, challengeID)
}

func TestCreateChallengeLogWithConfirmation_ErrorCheckingIfTaken(t *testing.T) {
	mockChallengeRepo := new(MockChallengeRepository)
	mockImpactRepo := new(MockImpactRepository)
	service := NewChallengeService(mockChallengeRepo, mockImpactRepo)

	challengeID := "challenge1"
	userID := "user1"
	log := challenges.ChallengeLog{
		ChallengeID: challengeID,
		UserID:      userID,
	}

	mockChallengeRepo.On("IsChallengeTaken", userID, challengeID).Return(false, errors.New("db error"))

	err := service.CreateChallengeLogWithConfirmation(log)

	assert.Error(t, err)
	assert.Equal(t, "db error", err.Error())
}

func TestCreateChallengeLogWithConfirmation_ErrorCreatingChallengeLog(t *testing.T) {
	mockChallengeRepo := new(MockChallengeRepository)
	mockImpactRepo := new(MockImpactRepository)
	service := NewChallengeService(mockChallengeRepo, mockImpactRepo)

	challengeID := "challenge1"
	userID := "user1"
	log := challenges.ChallengeLog{
		ChallengeID: challengeID,
		UserID:      userID,
	}

	mockChallengeRepo.On("IsChallengeTaken", userID, challengeID).Return(false, nil)
	mockChallengeRepo.On("CreateChallengeLog", mock.Anything).Return(errors.New("db error"))

	err := service.CreateChallengeLogWithConfirmation(log)

	assert.Error(t, err)
	assert.Equal(t, "db error", err.Error())
}

func TestCreateChallengeLogWithConfirmation_ErrorCreatingChallengeConfirmation(t *testing.T) {
	mockChallengeRepo := new(MockChallengeRepository)
	mockImpactRepo := new(MockImpactRepository)
	service := NewChallengeService(mockChallengeRepo, mockImpactRepo)

	challengeID := "challenge1"
	userID := "user1"
	log := challenges.ChallengeLog{
		ChallengeID: challengeID,
		UserID:      userID,
	}

	mockChallengeRepo.On("IsChallengeTaken", userID, challengeID).Return(false, nil)
	mockChallengeRepo.On("CreateChallengeLog", mock.Anything).Return(nil)
	mockChallengeRepo.On("GetTasksByChallengeID", challengeID).Return([]challenges.ChallengeTask{
		{ID: "task1"},
		{ID: "task2"},
	}, nil)
	mockChallengeRepo.On("CreateChallengeConfirmation", mock.Anything).Return(errors.New("db error")).Once()

	err := service.CreateChallengeLogWithConfirmation(log)

	assert.Error(t, err)
	assert.Equal(t, "db error", err.Error())
}

func TestCreateChallengeLogWithConfirmation_ErrorIncrementingCounts(t *testing.T) {
	mockChallengeRepo := new(MockChallengeRepository)
	mockImpactRepo := new(MockImpactRepository)
	service := NewChallengeService(mockChallengeRepo, mockImpactRepo)

	challengeID := "challenge1"
	userID := "user1"
	log := challenges.ChallengeLog{
		ChallengeID: challengeID,
		UserID:      userID,
	}

	mockChallengeRepo.On("IsChallengeTaken", userID, challengeID).Return(false, nil)
	mockChallengeRepo.On("CreateChallengeLog", mock.Anything).Return(nil)
	mockChallengeRepo.On("GetTasksByChallengeID", challengeID).Return([]challenges.ChallengeTask{
		{ID: "task1"},
		{ID: "task2"},
	}, nil)
	mockChallengeRepo.On("CreateChallengeConfirmation", mock.Anything).Return(nil).Twice()
	mockChallengeRepo.On("IncrementChallengeCounts", challengeID, 2, true).Return(errors.New("db error"))

	err := service.CreateChallengeLogWithConfirmation(log)

	assert.Error(t, err)
	assert.Equal(t, "db error", err.Error())
}

func TestUpdateChallengeConfirmationProgress_Success(t *testing.T) {
	mockChallengeRepo := new(MockChallengeRepository)
	mockImpactRepo := new(MockImpactRepository)
	service := NewChallengeService(mockChallengeRepo, mockImpactRepo)

	confirmationID := "confirmation1"
	userID := "user1"
	challengeImgURL := "http://example.com/image.png"

	mockChallengeRepo.On("GetChallengeConfirmationByID", confirmationID).Return(challenges.ChallengeConfirmation{
		ID:     confirmationID,
		UserID: userID,
	}, nil)

	mockChallengeRepo.On("UpdateChallengeConfirmation", mock.Anything).Return(nil)

	err := service.UpdateChallengeConfirmationProgress(confirmationID, challengeImgURL, userID)

	assert.NoError(t, err)
	mockChallengeRepo.AssertCalled(t, "GetChallengeConfirmationByID", confirmationID)
	mockChallengeRepo.AssertCalled(t, "UpdateChallengeConfirmation", mock.Anything)
}

func TestUpdateChallengeConfirmationProgress_Unauthorized(t *testing.T) {
	mockChallengeRepo := new(MockChallengeRepository)
	mockImpactRepo := new(MockImpactRepository)
	service := NewChallengeService(mockChallengeRepo, mockImpactRepo)

	confirmationID := "confirmation1"
	userID := "user1"

	mockChallengeRepo.On("GetChallengeConfirmationByID", confirmationID).Return(challenges.ChallengeConfirmation{
		ID:     confirmationID,
		UserID: "anotherUser",
	}, nil)

	err := service.UpdateChallengeConfirmationProgress(confirmationID, "http://example.com/image.png", userID)

	assert.Error(t, err)
	assert.Equal(t, constant.ErrUnauthorized, err)
	mockChallengeRepo.AssertCalled(t, "GetChallengeConfirmationByID", confirmationID)
}

func TestUpdateChallengeConfirmationProgress_ErrorOnGetChallengeConfirmation(t *testing.T) {
	mockChallengeRepo := new(MockChallengeRepository)
	mockImpactRepo := new(MockImpactRepository)
	service := NewChallengeService(mockChallengeRepo, mockImpactRepo)

	confirmationID := "confirmation1"
	userID := "user1"

	mockChallengeRepo.On("GetChallengeConfirmationByID", confirmationID).Return(challenges.ChallengeConfirmation{}, errors.New("failed to get confirmation"))

	err := service.UpdateChallengeConfirmationProgress(confirmationID, "http://example.com/image.png", userID)

	assert.Error(t, err)
	assert.EqualError(t, err, "failed to get confirmation")
	mockChallengeRepo.AssertCalled(t, "GetChallengeConfirmationByID", confirmationID)
}

func TestUpdateChallengeConfirmationProgress_ErrorOnUpdateChallengeConfirmation(t *testing.T) {
	mockChallengeRepo := new(MockChallengeRepository)
	mockImpactRepo := new(MockImpactRepository)
	service := NewChallengeService(mockChallengeRepo, mockImpactRepo)

	confirmationID := "confirmation1"
	userID := "user1"
	challengeImgURL := "http://example.com/image.png"

	mockChallengeRepo.On("GetChallengeConfirmationByID", confirmationID).Return(challenges.ChallengeConfirmation{
		ID:     confirmationID,
		UserID: userID,
	}, nil)

	mockChallengeRepo.On("UpdateChallengeConfirmation", mock.Anything).Return(errors.New("failed to update confirmation"))

	err := service.UpdateChallengeConfirmationProgress(confirmationID, challengeImgURL, userID)

	assert.Error(t, err)
	assert.EqualError(t, err, "failed to update confirmation")
	mockChallengeRepo.AssertCalled(t, "GetChallengeConfirmationByID", confirmationID)
	mockChallengeRepo.AssertCalled(t, "UpdateChallengeConfirmation", mock.Anything)
}

func TestCheckAndUpdateChallengeLogStatusByConfirmation_AllDone(t *testing.T) {
	mockChallengeRepo := new(MockChallengeRepository)
	mockImpactRepo := new(MockImpactRepository)
	service := NewChallengeService(mockChallengeRepo, mockImpactRepo)

	confirmationID := "confirmation1"
	userID := "user1"
	challengeID := "challenge1"

	mockChallengeRepo.On("GetChallengeConfirmationByID", confirmationID).Return(challenges.ChallengeConfirmation{
		ID:              confirmationID,
		ChallengeTaskID: "task1",
	}, nil)

	mockChallengeRepo.On("GetChallengeTaskByID", "task1").Return(challenges.ChallengeTask{
		ID:          "task1",
		ChallengeID: challengeID,
	}, nil)

	mockChallengeRepo.On("GetChallengeLogByChallengeIDAndUserID", challengeID, userID).Return(challenges.ChallengeLog{
		ID: challengeID,
	}, nil)

	mockChallengeRepo.On("GetConfirmationsByChallengeID", challengeID, userID).Return([]challenges.ChallengeConfirmation{
		{Status: "Done"},
		{Status: "Done"},
	}, nil)

	mockChallengeRepo.On("UpdateChallengeLog", mock.Anything).Return(nil)

	err := service.CheckAndUpdateChallengeLogStatusByConfirmation(confirmationID, userID)

	assert.NoError(t, err)
	mockChallengeRepo.AssertCalled(t, "GetConfirmationsByChallengeID", challengeID, userID)
	mockChallengeRepo.AssertCalled(t, "UpdateChallengeLog", mock.Anything)
}

func TestCheckAndUpdateChallengeLogStatusByConfirmation_ErrorOnGetChallengeTaskByID(t *testing.T) {
	mockChallengeRepo := new(MockChallengeRepository)
	mockImpactRepo := new(MockImpactRepository)
	service := NewChallengeService(mockChallengeRepo, mockImpactRepo)

	confirmationID := "confirmation1"
	userID := "user1"

	mockChallengeRepo.On("GetChallengeConfirmationByID", confirmationID).Return(challenges.ChallengeConfirmation{
		ID:              confirmationID,
		ChallengeTaskID: "task1",
	}, nil)

	mockChallengeRepo.On("GetChallengeTaskByID", "task1").Return(challenges.ChallengeTask{}, errors.New("task retrieval error"))

	err := service.CheckAndUpdateChallengeLogStatusByConfirmation(confirmationID, userID)

	assert.Error(t, err)
	assert.EqualError(t, err, "task retrieval error")
	mockChallengeRepo.AssertCalled(t, "GetChallengeTaskByID", "task1")
}

func TestCheckAndUpdateChallengeLogStatusByConfirmation_ErrorOnGetChallengeLog(t *testing.T) {
	mockChallengeRepo := new(MockChallengeRepository)
	mockImpactRepo := new(MockImpactRepository)
	service := NewChallengeService(mockChallengeRepo, mockImpactRepo)

	confirmationID := "confirmation1"
	userID := "user1"

	mockChallengeRepo.On("GetChallengeConfirmationByID", confirmationID).Return(challenges.ChallengeConfirmation{
		ID:              confirmationID,
		ChallengeTaskID: "task1",
	}, nil)

	mockChallengeRepo.On("GetChallengeTaskByID", "task1").Return(challenges.ChallengeTask{
		ID:          "task1",
		ChallengeID: "challenge1",
	}, nil)

	mockChallengeRepo.On("GetChallengeLogByChallengeIDAndUserID", "challenge1", userID).Return(challenges.ChallengeLog{}, errors.New("log retrieval error"))

	err := service.CheckAndUpdateChallengeLogStatusByConfirmation(confirmationID, userID)

	assert.Error(t, err)
	assert.EqualError(t, err, "log retrieval error")
	mockChallengeRepo.AssertCalled(t, "GetChallengeLogByChallengeIDAndUserID", "challenge1", userID)
}

func TestCheckAndUpdateChallengeLogStatusByConfirmation_ErrorOnGetConfirmations(t *testing.T) {
	mockChallengeRepo := new(MockChallengeRepository)
	mockImpactRepo := new(MockImpactRepository)
	service := NewChallengeService(mockChallengeRepo, mockImpactRepo)

	confirmationID := "confirmation1"
	userID := "user1"

	mockChallengeRepo.On("GetChallengeConfirmationByID", confirmationID).Return(challenges.ChallengeConfirmation{
		ID:              confirmationID,
		ChallengeTaskID: "task1",
	}, nil)

	mockChallengeRepo.On("GetChallengeTaskByID", "task1").Return(challenges.ChallengeTask{
		ID:          "task1",
		ChallengeID: "challenge1",
	}, nil)

	mockChallengeRepo.On("GetChallengeLogByChallengeIDAndUserID", "challenge1", userID).Return(challenges.ChallengeLog{
		ID: "log1",
	}, nil)

	mockChallengeRepo.On("GetConfirmationsByChallengeID", "challenge1", userID).Return([]challenges.ChallengeConfirmation{}, errors.New("confirmation retrieval error"))

	err := service.CheckAndUpdateChallengeLogStatusByConfirmation(confirmationID, userID)

	assert.Error(t, err)
	assert.EqualError(t, err, "confirmation retrieval error")
	mockChallengeRepo.AssertCalled(t, "GetConfirmationsByChallengeID", "challenge1", userID)
}

func TestCheckAndUpdateChallengeLogStatusByConfirmation_ErrorOnUpdateChallengeLog(t *testing.T) {
	mockChallengeRepo := new(MockChallengeRepository)
	mockImpactRepo := new(MockImpactRepository)
	service := NewChallengeService(mockChallengeRepo, mockImpactRepo)

	confirmationID := "confirmation1"
	userID := "user1"
	challengeID := "challenge1"

	mockChallengeRepo.On("GetChallengeConfirmationByID", confirmationID).Return(challenges.ChallengeConfirmation{
		ID:              confirmationID,
		ChallengeTaskID: "task1",
	}, nil)

	mockChallengeRepo.On("GetChallengeTaskByID", "task1").Return(challenges.ChallengeTask{
		ID:          "task1",
		ChallengeID: challengeID,
	}, nil)

	mockChallengeRepo.On("GetChallengeLogByChallengeIDAndUserID", challengeID, userID).Return(challenges.ChallengeLog{
		ID:     "log1",
		Status: "InProgress",
	}, nil)

	mockChallengeRepo.On("GetConfirmationsByChallengeID", challengeID, userID).Return([]challenges.ChallengeConfirmation{
		{Status: "Done"},
		{Status: "Done"},
	}, nil)

	mockChallengeRepo.On("UpdateChallengeLog", mock.Anything).Return(errors.New("update log error"))

	err := service.CheckAndUpdateChallengeLogStatusByConfirmation(confirmationID, userID)

	assert.Error(t, err)
	assert.EqualError(t, err, "update log error")
	mockChallengeRepo.AssertCalled(t, "UpdateChallengeLog", mock.Anything)
}

func TestCheckAndUpdateChallengeLogStatusByConfirmation_NotAllDone(t *testing.T) {
	mockChallengeRepo := new(MockChallengeRepository)
	mockImpactRepo := new(MockImpactRepository)
	service := NewChallengeService(mockChallengeRepo, mockImpactRepo)

	confirmationID := "confirmation1"
	userID := "user1"
	challengeID := "challenge1"

	mockChallengeRepo.On("GetChallengeConfirmationByID", confirmationID).Return(challenges.ChallengeConfirmation{
		ID:              confirmationID,
		ChallengeTaskID: "task1",
	}, nil)

	mockChallengeRepo.On("GetChallengeTaskByID", "task1").Return(challenges.ChallengeTask{
		ID:          "task1",
		ChallengeID: challengeID,
	}, nil)

	mockChallengeRepo.On("GetChallengeLogByChallengeIDAndUserID", challengeID, userID).Return(challenges.ChallengeLog{
		ID: challengeID,
	}, nil)

	mockChallengeRepo.On("GetConfirmationsByChallengeID", challengeID, userID).Return([]challenges.ChallengeConfirmation{
		{Status: "Done"},
		{Status: "Pending"},
	}, nil)

	err := service.CheckAndUpdateChallengeLogStatusByConfirmation(confirmationID, userID)

	assert.NoError(t, err)
	mockChallengeRepo.AssertCalled(t, "GetConfirmationsByChallengeID", challengeID, userID)
	mockChallengeRepo.AssertNotCalled(t, "UpdateChallengeLog", mock.Anything)
}

func TestCheckAndUpdateChallengeLogStatusByConfirmation_ErrorHandling(t *testing.T) {
	mockChallengeRepo := new(MockChallengeRepository)
	mockImpactRepo := new(MockImpactRepository)
	service := NewChallengeService(mockChallengeRepo, mockImpactRepo)

	confirmationID := "confirmation1"
	userID := "user1"

	mockChallengeRepo.On("GetChallengeConfirmationByID", confirmationID).Return(challenges.ChallengeConfirmation{}, errors.New("some error"))

	err := service.CheckAndUpdateChallengeLogStatusByConfirmation(confirmationID, userID)

	assert.Error(t, err)
	mockChallengeRepo.AssertCalled(t, "GetChallengeConfirmationByID", confirmationID)
}

func TestClaimRewards_Success(t *testing.T) {
	mockChallengeRepo := new(MockChallengeRepository)
	mockImpactRepo := new(MockImpactRepository)
	service := NewChallengeService(mockChallengeRepo, mockImpactRepo)

	challengeLogID := "log1"
	userID := "user1"
	challengeID := "challenge1"

	mockChallengeRepo.On("IsRewardClaimed", challengeLogID).Return(false, nil)
	mockChallengeRepo.On("GetChallengeIDByLogID", challengeLogID).Return(challengeID, nil)
	mockChallengeRepo.On("GetChallengeRewards", challengeID).Return(100, 50, nil)
	mockChallengeRepo.On("UpdateRewardsGiven", challengeLogID).Return(nil)
	mockChallengeRepo.On("AddUserRewards", userID, 100, 50).Return(nil)

	err := service.ClaimRewards(challengeLogID, userID)

	assert.NoError(t, err)
	mockChallengeRepo.AssertCalled(t, "IsRewardClaimed", challengeLogID)
	mockChallengeRepo.AssertCalled(t, "GetChallengeIDByLogID", challengeLogID)
	mockChallengeRepo.AssertCalled(t, "GetChallengeRewards", challengeID)
	mockChallengeRepo.AssertCalled(t, "UpdateRewardsGiven", challengeLogID)
	mockChallengeRepo.AssertCalled(t, "AddUserRewards", userID, 100, 50)
}

func TestClaimRewards_AlreadyClaimed(t *testing.T) {
	mockChallengeRepo := new(MockChallengeRepository)
	mockImpactRepo := new(MockImpactRepository)
	service := NewChallengeService(mockChallengeRepo, mockImpactRepo)

	challengeLogID := "log1"

	mockChallengeRepo.On("IsRewardClaimed", challengeLogID).Return(true, nil)

	err := service.ClaimRewards(challengeLogID, "user1")

	assert.Error(t, err)
	assert.Equal(t, constant.ErrRewardAlreadyClaimed, err)
	mockChallengeRepo.AssertCalled(t, "IsRewardClaimed", challengeLogID)
}

func TestClaimRewards_ErrorOnIsRewardClaimed(t *testing.T) {
	mockChallengeRepo := new(MockChallengeRepository)
	mockImpactRepo := new(MockImpactRepository)
	service := NewChallengeService(mockChallengeRepo, mockImpactRepo)

	challengeLogID := "log1"

	mockChallengeRepo.On("IsRewardClaimed", challengeLogID).Return(false, errors.New("database error"))

	err := service.ClaimRewards(challengeLogID, "user1")

	assert.Error(t, err)
	assert.EqualError(t, err, "database error")
	mockChallengeRepo.AssertCalled(t, "IsRewardClaimed", challengeLogID)
}

func TestClaimRewards_ErrorOnGetChallengeIDByLogID(t *testing.T) {
	mockChallengeRepo := new(MockChallengeRepository)
	mockImpactRepo := new(MockImpactRepository)
	service := NewChallengeService(mockChallengeRepo, mockImpactRepo)

	challengeLogID := "log1"

	mockChallengeRepo.On("IsRewardClaimed", challengeLogID).Return(false, nil)
	mockChallengeRepo.On("GetChallengeIDByLogID", challengeLogID).Return("", errors.New("database error"))

	err := service.ClaimRewards(challengeLogID, "user1")

	assert.Error(t, err)
	assert.EqualError(t, err, "database error")
	mockChallengeRepo.AssertCalled(t, "IsRewardClaimed", challengeLogID)
	mockChallengeRepo.AssertCalled(t, "GetChallengeIDByLogID", challengeLogID)
}

func TestClaimRewards_ErrorOnGetChallengeRewards(t *testing.T) {
	mockChallengeRepo := new(MockChallengeRepository)
	mockImpactRepo := new(MockImpactRepository)
	service := NewChallengeService(mockChallengeRepo, mockImpactRepo)

	challengeLogID := "log1"
	challengeID := "challenge1"

	mockChallengeRepo.On("IsRewardClaimed", challengeLogID).Return(false, nil)
	mockChallengeRepo.On("GetChallengeIDByLogID", challengeLogID).Return(challengeID, nil)
	mockChallengeRepo.On("GetChallengeRewards", challengeID).Return(0, 0, errors.New("database error"))

	err := service.ClaimRewards(challengeLogID, "user1")

	assert.Error(t, err)
	assert.EqualError(t, err, "database error")
	mockChallengeRepo.AssertCalled(t, "IsRewardClaimed", challengeLogID)
	mockChallengeRepo.AssertCalled(t, "GetChallengeIDByLogID", challengeLogID)
	mockChallengeRepo.AssertCalled(t, "GetChallengeRewards", challengeID)
}

func TestClaimRewards_ErrorOnUpdateRewardsGiven(t *testing.T) {
	mockChallengeRepo := new(MockChallengeRepository)
	mockImpactRepo := new(MockImpactRepository)
	service := NewChallengeService(mockChallengeRepo, mockImpactRepo)

	challengeLogID := "log1"
	challengeID := "challenge1"

	mockChallengeRepo.On("IsRewardClaimed", challengeLogID).Return(false, nil)
	mockChallengeRepo.On("GetChallengeIDByLogID", challengeLogID).Return(challengeID, nil)
	mockChallengeRepo.On("GetChallengeRewards", challengeID).Return(100, 50, nil)
	mockChallengeRepo.On("UpdateRewardsGiven", challengeLogID).Return(errors.New("database error"))

	err := service.ClaimRewards(challengeLogID, "user1")

	assert.Error(t, err)
	assert.EqualError(t, err, "database error")
	mockChallengeRepo.AssertCalled(t, "IsRewardClaimed", challengeLogID)
	mockChallengeRepo.AssertCalled(t, "GetChallengeIDByLogID", challengeLogID)
	mockChallengeRepo.AssertCalled(t, "GetChallengeRewards", challengeID)
	mockChallengeRepo.AssertCalled(t, "UpdateRewardsGiven", challengeLogID)
}

func TestClaimRewards_ErrorOnAddUserRewards(t *testing.T) {
	mockChallengeRepo := new(MockChallengeRepository)
	mockImpactRepo := new(MockImpactRepository)
	service := NewChallengeService(mockChallengeRepo, mockImpactRepo)

	challengeLogID := "log1"
	challengeID := "challenge1"

	mockChallengeRepo.On("IsRewardClaimed", challengeLogID).Return(false, nil)
	mockChallengeRepo.On("GetChallengeIDByLogID", challengeLogID).Return(challengeID, nil)
	mockChallengeRepo.On("GetChallengeRewards", challengeID).Return(100, 50, nil)
	mockChallengeRepo.On("UpdateRewardsGiven", challengeLogID).Return(nil)
	mockChallengeRepo.On("AddUserRewards", "user1", 100, 50).Return(errors.New("database error"))

	err := service.ClaimRewards(challengeLogID, "user1")

	assert.Error(t, err)
	assert.EqualError(t, err, "database error")
	mockChallengeRepo.AssertCalled(t, "IsRewardClaimed", challengeLogID)
	mockChallengeRepo.AssertCalled(t, "GetChallengeIDByLogID", challengeLogID)
	mockChallengeRepo.AssertCalled(t, "GetChallengeRewards", challengeID)
	mockChallengeRepo.AssertCalled(t, "UpdateRewardsGiven", challengeLogID)
	mockChallengeRepo.AssertCalled(t, "AddUserRewards", "user1", 100, 50)
}

func TestGetActiveChallenges_Success(t *testing.T) {
	mockChallengeRepo := new(MockChallengeRepository)
	service := NewChallengeService(mockChallengeRepo, nil)

	userID := "user1"
	page, perPage := 1, 10
	difficulty := "medium"
	title := "Save"

	mockChallengeRepo.On("GetChallengeLogByUserID", userID, page, perPage, difficulty, title).Return([]challenges.ChallengeLog{
		{ID: "log1", UserID: userID},
	}, 1, nil)

	challenges, total, err := service.GetActiveChallenges(userID, page, perPage, difficulty, title)

	assert.NoError(t, err)
	assert.Equal(t, 1, total)
	assert.Len(t, challenges, 1)

	mockChallengeRepo.AssertCalled(t, "GetChallengeLogByUserID", userID, page, perPage, difficulty, title)
}


func TestGetUnclaimedChallenges_Success(t *testing.T) {
	mockChallengeRepo := new(MockChallengeRepository)
	service := NewChallengeService(mockChallengeRepo, nil)

	userID := "user1"
	isAdmin := false
	page, limit := 1, 10
	difficulty := "easy"
	title := "Plant"

	mockChallengeRepo.On("GetUnclaimedChallenges", userID, isAdmin, page, limit, difficulty, title).Return([]challenges.Challenge{
		{ID: "challenge1", Title: "Plant Trees"},
	}, 1, nil)

	challenges, total, err := service.GetUnclaimedChallenges(userID, isAdmin, page, limit, difficulty, title)

	assert.NoError(t, err)
	assert.Equal(t, 1, total)
	assert.Len(t, challenges, 1)
	mockChallengeRepo.AssertCalled(t, "GetUnclaimedChallenges", userID, isAdmin, page, limit, difficulty, title)
}

func TestGetChallengeDetailsWithConfirmations_Success(t *testing.T) {
	mockChallengeRepo := new(MockChallengeRepository)
	service := NewChallengeService(mockChallengeRepo, nil)

	userID := "user1"
	challengeLogID := "log1"

	mockChallengeRepo.On("GetChallengeLogByID", challengeLogID).Return(challenges.ChallengeLog{
		ID:          challengeLogID,
		UserID:      userID,
		ChallengeID: "challenge1",
	}, nil)

	mockChallengeRepo.On("GetConfirmationsByChallengeID", "challenge1", userID).Return([]challenges.ChallengeConfirmation{
		{ID: "conf1", ChallengeTaskID: "task1", UserID: userID, Status: "Done"},
	}, nil)

	result, err := service.GetChallengeDetailsWithConfirmations(userID, challengeLogID)

	assert.NoError(t, err)
	assert.Equal(t, challengeLogID, result.ChallengeLog.ID)
	assert.Len(t, result.Confirmations, 1)
	mockChallengeRepo.AssertCalled(t, "GetChallengeLogByID", challengeLogID)
	mockChallengeRepo.AssertCalled(t, "GetConfirmationsByChallengeID", "challenge1", userID)
}

func TestGetChallengeDetailsWithConfirmations_Unauthorized(t *testing.T) {
	mockChallengeRepo := new(MockChallengeRepository)
	service := NewChallengeService(mockChallengeRepo, nil)

	userID := "user1"
	challengeLogID := "log1"

	mockChallengeRepo.On("GetChallengeLogByID", challengeLogID).Return(challenges.ChallengeLog{
		ID:     challengeLogID,
		UserID: "anotherUser",
	}, nil)

	result, err := service.GetChallengeDetailsWithConfirmations(userID, challengeLogID)

	assert.Error(t, err)
	assert.Equal(t, constant.ErrUnauthorized, err)
	assert.Empty(t, result)
	mockChallengeRepo.AssertCalled(t, "GetChallengeLogByID", challengeLogID)
}

func TestGetChallengeDetailsWithConfirmations_ErrorOnGetChallengeLog(t *testing.T) {
	mockChallengeRepo := new(MockChallengeRepository)
	service := NewChallengeService(mockChallengeRepo, nil)

	userID := "user1"
	challengeLogID := "log1"

	mockChallengeRepo.On("GetChallengeLogByID", challengeLogID).Return(challenges.ChallengeLog{}, errors.New("database error"))

	result, err := service.GetChallengeDetailsWithConfirmations(userID, challengeLogID)

	assert.Error(t, err)
	assert.EqualError(t, err, "database error")
	assert.Empty(t, result)
	mockChallengeRepo.AssertCalled(t, "GetChallengeLogByID", challengeLogID)
}

func TestGetChallengeDetailsWithConfirmations_ErrorOnGetConfirmations(t *testing.T) {
	mockChallengeRepo := new(MockChallengeRepository)
	service := NewChallengeService(mockChallengeRepo, nil)

	userID := "user1"
	challengeLogID := "log1"

	mockChallengeRepo.On("GetChallengeLogByID", challengeLogID).Return(challenges.ChallengeLog{
		ID:          challengeLogID,
		UserID:      userID,
		ChallengeID: "challenge1",
	}, nil)

	mockChallengeRepo.On("GetConfirmationsByChallengeID", "challenge1", userID).Return([]challenges.ChallengeConfirmation{}, errors.New("database error"))

	result, err := service.GetChallengeDetailsWithConfirmations(userID, challengeLogID)

	assert.Error(t, err)
	assert.EqualError(t, err, "database error")
	assert.Empty(t, result)
	mockChallengeRepo.AssertCalled(t, "GetChallengeLogByID", challengeLogID)
	mockChallengeRepo.AssertCalled(t, "GetConfirmationsByChallengeID", "challenge1", userID)
}

func TestGetChallengeDetails_Success(t *testing.T) {
	mockChallengeRepo := new(MockChallengeRepository)
	service := NewChallengeService(mockChallengeRepo, nil)

	challengeID := "challenge1"

	mockChallengeRepo.On("GetChallengeByID", challengeID).Return(challenges.Challenge{
			ID:               challengeID,
			Title:            "Plant Trees",
			Difficulty:       "medium",
			ChallengeImg:     "image_url",
			Description:      "Description of challenge",
			DurationDays:     30,
			Exp:              100,
			Coin:             50,
			ParticipantCount: 10,
			ActionCount:      3,
	}, nil)

	mockChallengeRepo.On("GetTasksByChallengeIDforUser", challengeID).Return([]challenges.ChallengeTask{
			{ID: "task1", Name: "Task 1", DayNumber: 1, TaskDescription: "Description for task 1"},
	}, nil)

	details, err := service.GetChallengeDetails(challengeID)

	assert.NoError(t, err)
	assert.Equal(t, challengeID, details.ID)
	assert.Equal(t, 10, details.ParticipantCount)
	assert.Equal(t, 3, details.ActionCount)
	assert.Len(t, details.Tasks, 1)
	mockChallengeRepo.AssertCalled(t, "GetChallengeByID", challengeID)
	mockChallengeRepo.AssertCalled(t, "GetTasksByChallengeIDforUser", challengeID)
}

func TestGetChallengeDetails_NoTasks(t *testing.T) {
	mockChallengeRepo := new(MockChallengeRepository)
	service := NewChallengeService(mockChallengeRepo, nil)

	challengeID := "challenge1"

	mockChallengeRepo.On("GetChallengeByID", challengeID).Return(challenges.Challenge{
		ID: "challenge1",
	}, nil)

	mockChallengeRepo.On("GetTasksByChallengeIDforUser", challengeID).Return([]challenges.ChallengeTask{}, nil)

	details, err := service.GetChallengeDetails(challengeID)

	assert.NoError(t, err)
	assert.Equal(t, challengeID, details.ID)
	assert.Len(t, details.Tasks, 0)
	mockChallengeRepo.AssertCalled(t, "GetChallengeByID", challengeID)
	mockChallengeRepo.AssertCalled(t, "GetTasksByChallengeIDforUser", challengeID)
}

func TestGetChallengeDetails_ErrorOnGetChallengeByID(t *testing.T) {
	mockChallengeRepo := new(MockChallengeRepository)
	service := NewChallengeService(mockChallengeRepo, nil)

	challengeID := "challenge1"

	mockChallengeRepo.On("GetChallengeByID", challengeID).Return(challenges.Challenge{}, errors.New("database error"))

	details, err := service.GetChallengeDetails(challengeID)

	assert.Error(t, err)
	assert.EqualError(t, err, "database error")
	assert.Empty(t, details)
	mockChallengeRepo.AssertCalled(t, "GetChallengeByID", challengeID)
}

func TestGetChallengeDetails_ErrorOnGetTasks(t *testing.T) {
	mockChallengeRepo := new(MockChallengeRepository)
	service := NewChallengeService(mockChallengeRepo, nil)

	challengeID := "challenge1"

	mockChallengeRepo.On("GetChallengeByID", challengeID).Return(challenges.Challenge{
		ID: challengeID,
	}, nil)

	mockChallengeRepo.On("GetTasksByChallengeIDforUser", challengeID).Return([]challenges.ChallengeTask{}, errors.New("database error"))

	details, err := service.GetChallengeDetails(challengeID)

	assert.Error(t, err)
	assert.EqualError(t, err, "database error")
	assert.Empty(t, details)
	mockChallengeRepo.AssertCalled(t, "GetChallengeByID", challengeID)
	mockChallengeRepo.AssertCalled(t, "GetTasksByChallengeIDforUser", challengeID)
}
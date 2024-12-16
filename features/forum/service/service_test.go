package service

import (
	"errors"
	"greenenvironment/constant"
	"greenenvironment/features/forum"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockForumRepository is a mock implementation of forum.ForumRepositoryInterface
type MockForumRepository struct {
	mock.Mock
}

func (m *MockForumRepository) GetAllForum() ([]forum.Forum, error) {
	args := m.Called()
	return args.Get(0).([]forum.Forum), args.Error(1)
}

func (m *MockForumRepository) GetAllByPage(page int) ([]forum.ForumGetAll, int, error) {
	args := m.Called(page)
	return args.Get(0).([]forum.ForumGetAll), args.Int(1), args.Error(2)
}

func (m *MockForumRepository) GetForumByID(ID string) (forum.Forum, error) {
	args := m.Called(ID)
	return args.Get(0).(forum.Forum), args.Error(1)
}

func (m *MockForumRepository) PostForum(f forum.Forum) error {
	args := m.Called(f)
	return args.Error(0)
}

func (m *MockForumRepository) UpdateForum(f forum.EditForum) error {
	args := m.Called(f)
	return args.Error(0)
}

func (m *MockForumRepository) DeleteForum(forumID string) error {
	args := m.Called(forumID)
	return args.Error(0)
}

func (m *MockForumRepository) PostMessageForum(mf forum.MessageForum) error {
	args := m.Called(mf)
	return args.Error(0)
}

func (m *MockForumRepository) DeleteMessageForum(messageID string) error {
	args := m.Called(messageID)
	return args.Error(0)
}

func (m *MockForumRepository) GetMessageForumByID(ID string) (forum.MessageForum, error) {
	args := m.Called(ID)
	return args.Get(0).(forum.MessageForum), args.Error(1)
}

func (m *MockForumRepository) UpdateMessageForum(md forum.EditMessage) error {
	args := m.Called(md)
	return args.Error(0)
}

func (m *MockForumRepository) GetForumByUserID(ID string, page int) ([]forum.Forum, int, error) {
	args := m.Called(ID, page)
	return args.Get(0).([]forum.Forum), args.Int(1), args.Error(2)
}

func (m *MockForumRepository) GetMessagesByForumID(forumID string) ([]forum.MessageForum, error) {
	args := m.Called(forumID)
	return args.Get(0).([]forum.MessageForum), args.Error(1)
}

func (m *MockForumRepository) GetMessagesByForumIDWithPagination(forumID string, page int, pageSize int) ([]forum.MessageForum, error) {
	args := m.Called(forumID, page, pageSize)
	return args.Get(0).([]forum.MessageForum), args.Error(1)
}

func TestGetAllForum(t *testing.T) {
	mockRepo := new(MockForumRepository)
	mockRepo.On("GetAllForum").Return([]forum.Forum{}, nil)

	service := NewForumService(mockRepo)
	forums, err := service.GetAllForum()

	assert.NoError(t, err)
	assert.NotNil(t, forums)
	mockRepo.AssertExpectations(t)
}

func TestGetForumByID(t *testing.T) {
	mockRepo := new(MockForumRepository)
	mockRepo.On("GetForumByID", "1").Return(forum.Forum{ID: "1"}, nil)

	service := NewForumService(mockRepo)
	forum, err := service.GetForumByID("1")

	assert.NoError(t, err)
	assert.Equal(t, "1", forum.ID)
	mockRepo.AssertExpectations(t)
}

func TestPostForum(t *testing.T) {
	mockRepo := new(MockForumRepository)
	mockForum := forum.Forum{Title: "Test", Description: "Test Description"}
	mockRepo.On("PostForum", mockForum).Return(nil)

	service := NewForumService(mockRepo)
	err := service.PostForum(mockForum)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUpdateForum(t *testing.T) {
	mockRepo := new(MockForumRepository)
	mockForum := forum.EditForum{ID: "1", Title: "Updated Title", Description: "Updated Description"}
	mockRepo.On("UpdateForum", mockForum).Return(nil)

	service := NewForumService(mockRepo)
	err := service.UpdateForum(mockForum)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteForum(t *testing.T) {
	mockRepo := new(MockForumRepository)
	mockRepo.On("DeleteForum", "1").Return(nil)

	service := NewForumService(mockRepo)
	err := service.DeleteForum("1")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestGetAllByPage(t *testing.T) {
	mockRepo := new(MockForumRepository)
	service := NewForumService(mockRepo)

	expectedForums := []forum.ForumGetAll{{ID: "1", Title: "Forum 1"}, {ID: "2", Title: "Forum 2"}}

	// Positive case
	mockRepo.On("GetAllByPage", 1).Return(expectedForums, 2, nil)

	result, totalPages, err := service.GetAllByPage(1)

	assert.NoError(t, err)
	assert.Equal(t, expectedForums, result)
	assert.Equal(t, 2, totalPages)
	mockRepo.AssertExpectations(t)
}

func TestPostMessageForum(t *testing.T) {
	mockRepo := new(MockForumRepository)
	service := NewForumService(mockRepo)

	validMessage := forum.MessageForum{ForumID: "1", Message: "This is a message"}
	invalidMessage := forum.MessageForum{ForumID: "", Message: ""}

	// Positive case
	mockRepo.On("PostMessageForum", validMessage).Return(nil)
	err := service.PostMessageForum(validMessage)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)

	// Negative case: invalid data
	err = service.PostMessageForum(invalidMessage)
	assert.Error(t, err)
	assert.Equal(t, constant.ErrFieldData, err)
}

func TestDeleteMessageForum(t *testing.T) {
	mockRepo := new(MockForumRepository)
	service := NewForumService(mockRepo)

	// Positive case
	mockRepo.On("DeleteMessageForum", "1").Return(nil)
	err := service.DeleteMessageForum("1")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)

	// Negative case: repo error
	mockRepo.On("DeleteMessageForum", "2").Return(errors.New("delete failed"))
	err = service.DeleteMessageForum("2")

	assert.Error(t, err)
	assert.EqualError(t, err, "delete failed")
}

func TestGetMessageForumByID(t *testing.T) {
	mockRepo := new(MockForumRepository)
	service := NewForumService(mockRepo)

	expectedMessage := forum.MessageForum{ID: "1", ForumID: "1", Message: "This is a message"}

	// Positive case
	mockRepo.On("GetMessageForumByID", "1").Return(expectedMessage, nil)

	result, err := service.GetMessageForumByID("1")

	assert.NoError(t, err)
	assert.Equal(t, expectedMessage, result)
	mockRepo.AssertExpectations(t)

	// Negative case: empty ID
	_, err = service.GetMessageForumByID("")
	assert.Error(t, err)
	assert.Equal(t, constant.ErrFieldData, err)
}

func TestUpdateMessageForum(t *testing.T) {
	mockRepo := new(MockForumRepository)
	service := NewForumService(mockRepo)

	validMessage := forum.EditMessage{ID: "1", Message: "Updated message"}
	invalidMessage := forum.EditMessage{ID: "", Message: ""}

	// Positive case
	mockRepo.On("UpdateMessageForum", validMessage).Return(nil)
	err := service.UpdateMessageForum(validMessage)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)

	// Negative case: invalid data
	err = service.UpdateMessageForum(invalidMessage)
	assert.Error(t, err)
	assert.Equal(t, constant.ErrFieldData, err)
}

func TestGetForumByUserID(t *testing.T) {
	mockRepo := new(MockForumRepository)
	service := NewForumService(mockRepo)

	expectedForums := []forum.Forum{{ID: "1", Title: "User Forum 1"}, {ID: "2", Title: "User Forum 2"}}

	// Positive case
	mockRepo.On("GetForumByUserID", "123", 1).Return(expectedForums, 2, nil)

	result, totalPages, err := service.GetForumByUserID("123", 1)

	assert.NoError(t, err)
	assert.Equal(t, expectedForums, result)
	assert.Equal(t, 2, totalPages)
	mockRepo.AssertExpectations(t)
}

func TestGetMessagesByForumID(t *testing.T) {
	mockRepo := new(MockForumRepository)
	service := NewForumService(mockRepo)
	mockRepo.On("GetMessagesByForumID", "123").Return([]forum.MessageForum{
		{ID: "1", Message: "Message 1"},
		{ID: "2", Message: "Message 2"},
	}, nil)

	messages, err := service.GetMessagesByForumID("123")

	assert.NoError(t, err)
	assert.Len(t, messages, 2)
	assert.Equal(t, "1", messages[0].ID)
	assert.Equal(t, "Message 1", messages[0].Message)

	mockRepo.AssertExpectations(t)
}

func TestGetMessagesByForumIDWithPagination(t *testing.T) {
	mockRepo := new(MockForumRepository)
	service := NewForumService(mockRepo)

	mockRepo.On("GetMessagesByForumIDWithPagination", "123", 1, 10).Return([]forum.MessageForum{
		{ID: "1", Message: "Message 1"},
		{ID: "2", Message: "Message 2"},
	}, nil)

	messages, err := service.GetMessagesByForumIDWithPagination("123", 1, 10)

	assert.NoError(t, err)
	assert.Len(t, messages, 2)
	assert.Equal(t, "1", messages[0].ID)
	assert.Equal(t, "Message 1", messages[0].Message)

	mockRepo.AssertExpectations(t)
}

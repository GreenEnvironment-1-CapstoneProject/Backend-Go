package service

import (
	"errors"
	"greenenvironment/features/leaderboard"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockLeaderboardRepository struct {
	mock.Mock
}

func (m *MockLeaderboardRepository) GetLeaderboard() ([]leaderboard.LeaderboardUser, error) {
	args := m.Called()

	if args.Get(0) != nil {
		return args.Get(0).([]leaderboard.LeaderboardUser), args.Error(1)
	}
	return nil, args.Error(1)
}

func TestGetLeaderboard(t *testing.T) {
	mockRepo := new(MockLeaderboardRepository)
	expectedData := []leaderboard.LeaderboardUser{
		{
			Rank:      1,
			ID:        "1",
			Name:      "John Doe",
			AvatarURL: "http://example.com/avatar1.png",
			Exp:       100,
		},
		{
			Rank:      2,
			ID:        "2",
			Name:      "Jane Smith",
			AvatarURL: "http://example.com/avatar2.png",
			Exp:       80,
		},
	}

	t.Run("success", func(t *testing.T) {
		mockRepo.On("GetLeaderboard").Return(expectedData, nil).Once()

		leaderboardService := NewLeaderboardService(mockRepo)
		result, err := leaderboardService.GetLeaderboard()

		assert.NoError(t, err)
		assert.Equal(t, expectedData, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockError := errors.New("database error")
		mockRepo.On("GetLeaderboard").Return(nil, mockError).Once()

		leaderboardService := NewLeaderboardService(mockRepo)
		result, err := leaderboardService.GetLeaderboard()

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, mockError, err)
		mockRepo.AssertExpectations(t)
	})
}

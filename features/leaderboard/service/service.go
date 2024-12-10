package service

import (
	"greenenvironment/features/leaderboard"
)

type LeaderboardService struct {
	repo leaderboard.LeaderboardRepositoryInterface
}

func NewLeaderboardService(repo leaderboard.LeaderboardRepositoryInterface) leaderboard.LeaderboardServiceInterface {
	return &LeaderboardService{repo: repo}
}

func (ls *LeaderboardService) GetLeaderboard() ([]leaderboard.LeaderboardUser, error) {
	return ls.repo.GetLeaderboard()
}

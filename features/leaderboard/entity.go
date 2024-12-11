package leaderboard

import (
	"github.com/labstack/echo/v4"
)

type LeaderboardUser struct {
	Rank      int
	ID        string
	Name      string
	AvatarURL string
	Exp       int
}

type LeaderboardRepositoryInterface interface {
	GetLeaderboard() ([]LeaderboardUser, error)
}

type LeaderboardServiceInterface interface {
	GetLeaderboard() ([]LeaderboardUser, error)
}

type LeaderboardControllerInterface interface {
	GetLeaderboard(c echo.Context) error
}

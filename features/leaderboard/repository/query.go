package repository

import (
	"greenenvironment/features/leaderboard"

	"gorm.io/gorm"
)

type LeaderboardData struct {
	DB *gorm.DB
}

func NewLeaderboardRepository(db *gorm.DB) leaderboard.LeaderboardRepositoryInterface {
	return &LeaderboardData{DB: db}
}

func (ld *LeaderboardData) GetLeaderboard() ([]leaderboard.LeaderboardUser, error) {
	var leaderboardData []leaderboard.LeaderboardUser
	query := `
			SELECT ROW_NUMBER() OVER (ORDER BY exp DESC, name ASC) AS ` + "`rank`" + `, id, name, avatar_url, exp
			FROM users
			LIMIT 10
	`
	if err := ld.DB.Raw(query).Scan(&leaderboardData).Error; err != nil {
		return nil, err
	}

	return leaderboardData, nil
}

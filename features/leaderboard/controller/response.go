package controller

import "greenenvironment/features/leaderboard"

type LeaderboardResponse struct {
	Rank int    `json:"rank"`
	ID   string `json:"id"`
	Name string `json:"name"`
	Exp  int    `json:"exp"`
}


func (l LeaderboardResponse) FromEntity(entities []leaderboard.LeaderboardUser) []LeaderboardResponse {
	var responses []LeaderboardResponse
	for _, entity := range entities {
		responses = append(responses, LeaderboardResponse{
			Rank: entity.Rank,
			ID:   entity.ID,
			Name: entity.Name,
			Exp:  entity.Exp,
		})
	}
	return responses
}

package controller

import "greenenvironment/features/challenges"

type ChallengeResponse struct {
	ID               string                      `json:"id"`
	Author           string                      `json:"author"`
	Title            string                      `json:"title"`
	Difficulty       string                      `json:"difficulty"`
	ChallengeImg     string                      `json:"challenge_img"`
	Description      string                      `json:"description"`
	DurationDays     int                         `json:"duration_days"`
	Exp              int                         `json:"exp"`
	Coin             int                         `json:"coin"`
	ImpactCategories []ChallengeImpactCategories `json:"categories"`
}

type ChallengeImpactCategories struct {
	ImpactCategory ImpactCategory `json:"impact_category"`
}

type ImpactCategory struct {
	Name        string `json:"name"`
	ImpactPoint int    `json:"impact_point"`
	Description string `json:"description"`
}

type MetadataResponse struct {
	CurrentPage int `json:"current_page"`
	TotalPage   int `json:"total_page"`
}

func (cr ChallengeResponse) ToResponse(challenge challenges.Challenge) ChallengeResponse {
	impactCategories := make([]ChallengeImpactCategories, len(challenge.ImpactCategories))
	for i, impactCategory := range challenge.ImpactCategories {
		impactCategories[i] = ChallengeImpactCategories{
			ImpactCategory: ImpactCategory{
				Name:        impactCategory.ImpactCategory.Name,
				ImpactPoint: impactCategory.ImpactCategory.ImpactPoint,
				Description: impactCategory.ImpactCategory.Description,
			},
		}
	}

	return ChallengeResponse{
		ID:               challenge.ID,
		Author:           challenge.Author,
		Title:            challenge.Title,
		Difficulty:       challenge.Difficulty,
		ChallengeImg:     challenge.ChallengeImg,
		Description:      challenge.Description,
		DurationDays:     challenge.DurationDays,
		Exp:              challenge.Exp,
		Coin:             challenge.Coin,
		ImpactCategories: impactCategories,
	}
}

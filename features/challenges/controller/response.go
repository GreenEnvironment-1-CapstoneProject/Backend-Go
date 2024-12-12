package controller

import (
	"greenenvironment/features/challenges"
	"time"
)

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
	DeletedAt        *string                     `json:"deleted_at"`
}

type ChallengeLogResponse struct {
	ID                      string                        `json:"id"`
	Challenge               ChallengeResponse             `json:"challenge"`
	UserID                  string                        `json:"user_id"`
	RewardsGiven            bool                          `json:"rewards_given"`
	Status                  string                        `json:"status"`
	StartDate               time.Time                     `json:"start_date"`
	Feed                    string                        `json:"feed"`
	ChallengeConfirmation   []ChallengeConfirmationResponse `json:"challenge_confirmation"`
}

type ChallengeConfirmationResponse struct {
	ID              string                `json:"id"`
	ChallengeTask   ChallengeTaskResponse `json:"challenge_task"`
	UserID          string                `json:"user_id"`
	Status          string                `json:"status"`
	SubmissionDate  time.Time             `json:"submission_date"`
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

	var deletedAt *string
	if challenge.DeletedAt != nil {
		formatted := challenge.DeletedAt.Format("2006-01-02 15:04:05")
		deletedAt = &formatted
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
		DeletedAt:        deletedAt,
	}
}

type ChallengeTaskResponse struct {
	ID              string `json:"id"`
	ChallengeID     string `json:"challenge_id"`
	Name            string `json:"name"`
	DayNumber       int    `json:"day_number"`
	TaskDescription string `json:"task_description"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}

func (ctr ChallengeTaskResponse) FromEntity(task challenges.ChallengeTask) ChallengeTaskResponse {
	return ChallengeTaskResponse{
		ID:              task.ID,
		ChallengeID:     task.ChallengeID,
		Name:            task.Name,
		DayNumber:       task.DayNumber,
		TaskDescription: task.TaskDescription,
		CreatedAt:       task.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:       task.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

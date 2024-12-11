package controller

type ChallengeRequest struct {
	Title            string   `form:"title" validate:"required"`
	Difficulty       string   `form:"difficulty" validate:"required"`
	Description      string   `form:"description" validate:"required"`
	DurationDays     int      `form:"duration_days" validate:"required"`
	Exp              int      `form:"exp" validate:"required"`
	Coin             int      `form:"coin" validate:"required"`
	ImpactCategories []string `form:"category_impact" validate:"required"`
}

type ChallengeTaskRequest struct {
	ChallengeID     string `json:"challenge_id" validate:"required"`
	Name            string `json:"name" validate:"required"`
	DayNumber       int    `json:"day_number" validate:"required,min=1"`
	TaskDescription string `json:"task_description" validate:"required"`
}

type ChallengeLogRequest struct {
	ChallengeID string `json:"challenge_id" validate:"required"`
	Feed        string `json:"feed"`
}

type ChallengeConfirmationRequest struct {
	ChallengeConfirmationID string `form:"challenge_confirmation_id" validate:"required"`
}

type ClaimRewardsRequest struct {
	ChallengeLogID string `json:"challenge_log_id" validate:"required"`
}

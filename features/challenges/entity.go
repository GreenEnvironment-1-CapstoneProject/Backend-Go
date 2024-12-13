package challenges

import (
	impactcategory "greenenvironment/features/impacts"
	"time"

	"github.com/labstack/echo/v4"
)

type Challenge struct {
	ID               string
	Author           string
	Title            string
	Difficulty       string
	ChallengeImg     string
	Description      string
	DurationDays     int
	Exp              int
	Coin             int
	ActionCount      int
	ParticipantCount int
	CreatedAt        time.Time
	UpdatedAt        time.Time
	ImpactCategories []ChallengeImpactCategory
	DeletedAt        *time.Time
}

type ChallengeImpactCategory struct {
	ID               string
	ChallengeID      string
	ImpactCategoryID string
	ImpactCategory   impactcategory.ImpactCategory
}

type ImpactCategory struct {
	ID          string
	Name        string
	ImpactPoint int
	Description string
}

type ChallengeTask struct {
	ID              string
	ChallengeID     string
	Name            string
	DayNumber       int
	TaskDescription string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type ChallengeLog struct {
	ID           string
	ChallengeID  string
	UserID       string
	Status       string
	StartDate    time.Time
	Feed         string
	RewardsGiven bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Challenge    Challenge
}

type ChallengeConfirmation struct {
	ID              string
	ChallengeTaskID string
	UserID          string
	Status          string
	ChallengeImg    string
	SubmissionDate  time.Time
	ChallengeTask   ChallengeTask
}

type ChallengeLogDetails struct {
	ChallengeLog  ChallengeLog
	Confirmations []ChallengeConfirmation
}

type ChallengeDetails struct {
	ID           string
	Title        string
	Difficulty   string
	ChallengeImg string
	Description  string
	DurationDays int
	Exp          int
	Coin         int
	Tasks        []ChallengeTask
}

type ChallengeWithCounts struct {
	Challenge
	ActionCount      int
	ParticipantCount int
}

type ChallengeRepoInterface interface {
	Create(Challenge) error
	GetAllByPage(page int) ([]Challenge, int, error)
	GetByID(id string) (Challenge, error)
	Update(Challenge) error
	Delete(id string) error

	CreateTask(task ChallengeTask) error
	GetTasksByChallengeID(challengeID string) ([]ChallengeTask, error)
	GetTaskByID(taskID string) (ChallengeTask, error)
	UpdateTask(task ChallengeTask) error
	DeleteTask(taskID string) error

	// User
	CreateChallengeLog(ChallengeLog) error
	CreateChallengeConfirmation(ChallengeConfirmation) error
	IsChallengeTaken(userID, challengeID string) (bool, error)
	IncrementChallengeCounts(challengeID string, actionCount int, participantIncrement bool) error
	GetChallengeConfirmationByID(confirmationID string) (ChallengeConfirmation, error)
	UpdateChallengeConfirmation(ChallengeConfirmation) error
	GetChallengeTaskByID(taskID string) (ChallengeTask, error)
	GetChallengeLogByChallengeIDAndUserID(challengeID, userID string) (ChallengeLog, error)
	GetConfirmationsByChallengeID(challengeID, userID string) ([]ChallengeConfirmation, error)
	UpdateChallengeLog(log ChallengeLog) error
	UpdateTaskAndChallengeStatus() error

	IsRewardClaimed(challengeLogID string) (bool, error)
	UpdateRewardsGiven(challengeLogID string) error
	AddUserRewards(userID string, exp int, coin int) error
	GetChallengeIDByLogID(challengeLogID string) (string, error)
	GetChallengeRewards(challengeID string) (int, int, error)

	GetChallengeLogByUserID(userID string, page, perPage int, difficulty, title string) ([]ChallengeLog, int, error)
	GetUnclaimedChallenges(userID string, isAdmin bool, page int, limit int, difficulty, title string) ([]Challenge, int, error)
	GetChallengeLogByID(challengeLogID string) (ChallengeLog, error)
	GetChallengeByID(challengeID string) (Challenge, error)
	GetTasksByChallengeIDforUser(challengeID string) ([]ChallengeTask, error)
}

type ChallengeServiceInterface interface {
	Create(Challenge) error
	GetAllByPage(page int) ([]Challenge, int, error)
	GetByID(id string) (Challenge, error)
	Update(Challenge) error
	Delete(challengeID string) error

	CreateTask(challengeID, name string, dayNumber int, taskDescription string) error
	GetAllTasksByChallengeID(challengeID string) ([]ChallengeTask, error)
	GetTaskByID(taskID string) (ChallengeTask, error)
	UpdateTask(taskID string, taskDescription string) error
	DeleteTask(taskID string) error

	// User
	CreateChallengeLogWithConfirmation(ChallengeLog) error
	UpdateChallengeConfirmationProgress(confirmationID, challengeImgURL, userID string) error
	CheckAndUpdateChallengeLogStatusByConfirmation(confirmationID, userID string) error
	ClaimRewards(challengeLogID, userID string) error

	GetActiveChallenges(userID string, page, perPage int, difficulty, title string) ([]ChallengeLog, int, error)
	GetUnclaimedChallenges(userID string, isAdmin bool, page, limit int, difficulty, title string) ([]Challenge, int, error)
	GetChallengeDetailsWithConfirmations(userID, challengeLogID string) (ChallengeLogDetails, error)
	GetChallengeDetails(challengeID string) (ChallengeDetails, error)
}

type ChallengeControllerInterface interface {
	Create(c echo.Context) error
	GetAll(c echo.Context) error
	GetByID(c echo.Context) error
	Update(c echo.Context) error
	Delete(c echo.Context) error

	CreateTask(c echo.Context) error
	GetAllTasksByChallengeID(c echo.Context) error
	GetTaskByID(c echo.Context) error
	UpdateTask(c echo.Context) error
	DeleteTask(c echo.Context) error

	// User
	CreateChallengeLog(c echo.Context) error
	UpdateChallengeConfirmationProgress(c echo.Context) error
	ClaimRewards(c echo.Context) error

	GetActiveChallenges(c echo.Context) error
	GetUnclaimedChallenges(c echo.Context) error
	GetChallengeDetailsWithConfirmations(c echo.Context) error
	GetChallengeDetails(c echo.Context) error
}

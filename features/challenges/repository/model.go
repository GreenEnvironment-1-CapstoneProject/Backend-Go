package repository

import (
	impactcategory "greenenvironment/features/impacts/repository"
	"time"

	"gorm.io/gorm"
)

type Challenge struct {
	*gorm.Model
	ID               string `gorm:"primaryKey;type:varchar(50);not null;column:id"`
	Author           string `gorm:"type:varchar(255);not null;column:author"`
	Title            string `gorm:"type:varchar(255);not null;column:title"`
	Difficulty       string `gorm:"type:varchar(255);not null;column:difficulty"`
	ChallengeImg     string `gorm:"type:varchar(255);not null;column:challenge_img"`
	Description      string `gorm:"type:varchar(255);not null;column:description"`
	DurationDays     int    `gorm:"type:int;not null;column:duration_days"`
	Exp              int    `gorm:"type:int;not null;column:exp"`
	Coin             int    `gorm:"type:int;not null;column:coin"`
	ActionCount      int    `gorm:"type:int;not null;default:0;column:action_count"`
	ParticipantCount int    `gorm:"type:int;not null;default:0;column:participant_count"`
	ImpactCategories []ChallengeImpactCategory `gorm:"foreignKey:ChallengeID;references:ID"`
}

type ChallengeImpactCategory struct {
	*gorm.Model
	ID               string                        `gorm:"primaryKey;type:varchar(50);not null;column:id"`
	ChallengeID      string                        `gorm:"type:varchar(50);not null;column:challenge_id"`
	ImpactCategoryID string                        `gorm:"type:varchar(50);not null;column:impact_category_id"`
	Challenge        Challenge                     `gorm:"foreignKey:ChallengeID;references:ID"`
	ImpactCategory   impactcategory.ImpactCategory `gorm:"foreignKey:ImpactCategoryID;references:ID"`
}

type ChallengeTask struct {
	*gorm.Model
	ID              string    `gorm:"primaryKey;type:varchar(50);not null;column:id"`
	ChallengeID     string    `gorm:"type:varchar(50);not null;column:challenge_id"`
	Name            string    `gorm:"type:varchar(255);not null;column:name"`
	DayNumber       int       `gorm:"type:int;not null;column:day_number"`
	TaskDescription string    `gorm:"type:text;not null;column:task_description"`
	Challenge       Challenge `gorm:"foreignKey:ChallengeID;references:ID"`
}

type ChallengeLog struct {
	*gorm.Model
	ID           string    `gorm:"primaryKey;type:varchar(50);not null;column:id"`
	ChallengeID  string    `gorm:"type:varchar(50);not null;column:challenge_id"`
	UserID       string    `gorm:"type:varchar(100);not null;column:user_id"`
	Status       string    `gorm:"type:enum('Progress','Done','Failed');not null;column:status"`
	StartDate    time.Time `gorm:"type:datetime;not null;column:start_date"`
	Feed         string    `gorm:"type:text;column:feed"`
	RewardsGiven bool      `gorm:"type:boolean;default:false;column:rewards_given"`
	Challenge    Challenge `gorm:"foreignKey:ChallengeID;references:ID"`
}

type ChallengeConfirmation struct {
	*gorm.Model
	ID              string        `gorm:"primaryKey;type:varchar(50);not null;column:id"`
	ChallengeTaskID string        `gorm:"type:varchar(50);not null;column:challenge_task_id"`
	UserID          string        `gorm:"type:varchar(50);not null;column:user_id"`
	Status          string        `gorm:"type:enum('Progress', 'Done', 'Failed');not null;column:status"`
	ChallengeImg    string        `gorm:"type:varchar(255);not null;column:challenge_img"`
	SubmissionDate  time.Time     `gorm:"type:datetime;not null;column:submission_date"`
	ChallengeTask   ChallengeTask `gorm:"foreignKey:ChallengeTaskID;references:ID"`
}

func (Challenge) TableName() string {
	return "challenges"
}

func (ChallengeImpactCategory) TableName() string {
	return "challenge_impact_categories"
}

func (ChallengeTask) TableName() string {
	return "challenge_tasks"
}

func (ChallengeLog) TableName() string {
	return "challenge_logs"
}

func (ChallengeConfirmation) TableName() string {
	return "challenge_confirmations"
}

package service

import (
	"greenenvironment/constant"
	"greenenvironment/features/challenges"
	"greenenvironment/features/impacts"
	"time"

	"github.com/google/uuid"
)

type ChallengeService struct {
	challengeRepo challenges.ChallengeRepoInterface
	impactRepo    impacts.ImpactRepositoryInterface
}

func NewChallengeService(cr challenges.ChallengeRepoInterface, ir impacts.ImpactRepositoryInterface) challenges.ChallengeServiceInterface {
	return &ChallengeService{
		challengeRepo: cr,
		impactRepo:    ir,
	}
}

func (cs *ChallengeService) Create(challenge challenges.Challenge) error {
	challenge.ID = uuid.New().String()
	for i, impact := range challenge.ImpactCategories {
		data, _ := cs.impactRepo.GetByID(impact.ImpactCategoryID)
		if data.ID == "" {
			return constant.ErrCreateProduct
		}
		impact.ID = uuid.New().String()
		challenge.ImpactCategories[i] = impact
	}

	err := cs.challengeRepo.Create(challenge)
	if err != nil {
		return err
	}
	return nil
}

func (cs *ChallengeService) GetAllByPage(page int) ([]challenges.Challenge, int, error) {
	challenges, totalPages, err := cs.challengeRepo.GetAllByPage(page)
	if err != nil {
		return nil, 0, err
	}
	return challenges, totalPages, nil
}

func (cs *ChallengeService) GetByID(id string) (challenges.Challenge, error) {
	return cs.challengeRepo.GetByID(id)
}

func (cs *ChallengeService) Update(challenge challenges.Challenge) error {
	for i, impact := range challenge.ImpactCategories {
		data, _ := cs.impactRepo.GetByID(impact.ImpactCategoryID)
		if data.ID == "" {
			return constant.ErrUpdateChallenge
		}
		impact.ID = uuid.New().String()
		challenge.ImpactCategories[i] = impact
	}

	existingChallenge, err := cs.challengeRepo.GetByID(challenge.ID)
	if err != nil {
		return err
	}

	challenge.Author = existingChallenge.Author

	if challenge.ChallengeImg == "" {
		challenge.ChallengeImg = existingChallenge.ChallengeImg
	}

	return cs.challengeRepo.Update(challenge)
}

func (cs *ChallengeService) Delete(challengeID string) error {
	_, err := cs.challengeRepo.GetByID(challengeID)
	if err != nil {
		return constant.ErrChallengeNotFound
	}

	return cs.challengeRepo.Delete(challengeID)
}

func (cs *ChallengeService) CreateTask(challengeID, name string, dayNumber int, taskDescription string) error {
	challenge, err := cs.challengeRepo.GetByID(challengeID)
	if err != nil || challenge.ID == "" {
		return constant.ErrChallengeNotFound
	}

	if dayNumber < 1 || dayNumber > challenge.DurationDays {
		return constant.ErrInvalidDayNumber
	}

	tasks, err := cs.challengeRepo.GetTasksByChallengeID(challengeID)
	if err != nil {
		return err
	}

	for _, task := range tasks {
		if task.DayNumber == dayNumber {
			return constant.ErrTaskAlreadyExists
		}
	}

	task := challenges.ChallengeTask{
		ID:              uuid.New().String(),
		ChallengeID:     challengeID,
		Name:            name,
		DayNumber:       dayNumber,
		TaskDescription: taskDescription,
	}

	return cs.challengeRepo.CreateTask(task)
}

func (cs *ChallengeService) GetAllTasksByChallengeID(challengeID string) ([]challenges.ChallengeTask, error) {
	return cs.challengeRepo.GetTasksByChallengeID(challengeID)
}

func (cs *ChallengeService) GetTaskByID(taskID string) (challenges.ChallengeTask, error) {
	task, err := cs.challengeRepo.GetTaskByID(taskID)
	if err != nil || task.ID == "" {
		return challenges.ChallengeTask{}, constant.ErrTaskNotFound
	}
	return task, nil
}

func (cs *ChallengeService) UpdateTask(taskID string, taskDescription string) error {
	task, err := cs.challengeRepo.GetTaskByID(taskID)
	if err != nil || task.ID == "" {
		return constant.ErrTaskNotFound
	}

	task.TaskDescription = taskDescription
	return cs.challengeRepo.UpdateTask(task)
}

func (cs *ChallengeService) DeleteTask(taskID string) error {
	task, err := cs.challengeRepo.GetTaskByID(taskID)
	if err != nil || task.ID == "" {
		return constant.ErrTaskNotFound
	}
	return cs.challengeRepo.DeleteTask(taskID)
}

// User
func (cs *ChallengeService) CreateChallengeLogWithConfirmation(log challenges.ChallengeLog) error {
	taken, err := cs.challengeRepo.IsChallengeTaken(log.UserID, log.ChallengeID)
	if err != nil {
		return err
	}
	if taken {
		return constant.ErrChallengeAlreadyTaken
	}

	log.ID = uuid.New().String()
	log.RewardsGiven = false

	err = cs.challengeRepo.CreateChallengeLog(log)
	if err != nil {
		return err
	}

	tasks, err := cs.challengeRepo.GetTasksByChallengeID(log.ChallengeID)
	if err != nil {
		return err
	}

	actionCount := len(tasks)
	for _, task := range tasks {
		confirmation := challenges.ChallengeConfirmation{
			ID:              uuid.New().String(),
			ChallengeTaskID: task.ID,
			UserID:          log.UserID,
			Status:          "Progress",
			ChallengeImg:    "",
			SubmissionDate:  time.Now(),
		}
		err := cs.challengeRepo.CreateChallengeConfirmation(confirmation)
		if err != nil {
			return err
		}
	}

	// Increment action_count and participant_count
	err = cs.challengeRepo.IncrementChallengeCounts(log.ChallengeID, actionCount, true)
	if err != nil {
		return err
	}

	return nil
}

func (cs *ChallengeService) UpdateChallengeConfirmationProgress(confirmationID, challengeImgURL, userID string) error {
	confirmation, err := cs.challengeRepo.GetChallengeConfirmationByID(confirmationID)
	if err != nil {
		return err
	}

	if confirmation.UserID != userID {
		return constant.ErrUnauthorized
	}

	confirmation.Status = "Done"
	confirmation.ChallengeImg = challengeImgURL
	confirmation.SubmissionDate = time.Now()

	err = cs.challengeRepo.UpdateChallengeConfirmation(confirmation)
	if err != nil {
		return err
	}

	return nil
}

func (cs *ChallengeService) CheckAndUpdateChallengeLogStatusByConfirmation(confirmationID, userID string) error {
	confirmation, err := cs.challengeRepo.GetChallengeConfirmationByID(confirmationID)
	if err != nil {
		return err
	}

	task, err := cs.challengeRepo.GetChallengeTaskByID(confirmation.ChallengeTaskID)
	if err != nil {
		return err
	}

	challengeLog, err := cs.challengeRepo.GetChallengeLogByChallengeIDAndUserID(task.ChallengeID, userID)
	if err != nil {
		return err
	}

	confirmations, err := cs.challengeRepo.GetConfirmationsByChallengeID(task.ChallengeID, userID)
	if err != nil {
		return err
	}

	allDone := true
	for _, conf := range confirmations {
		if conf.Status != "Done" {
			allDone = false
			break
		}
	}

	if allDone {
		challengeLog.Status = "Done"
		err := cs.challengeRepo.UpdateChallengeLog(challengeLog)
		if err != nil {
			return err
		}
	}

	return nil
}

func (cs *ChallengeService) ClaimRewards(challengeLogID, userID string) error {
	claimed, err := cs.challengeRepo.IsRewardClaimed(challengeLogID)
	if err != nil {
		return err
	}

	if claimed {
		return constant.ErrRewardAlreadyClaimed
	}

	challengeID, err := cs.challengeRepo.GetChallengeIDByLogID(challengeLogID)
	if err != nil {
		return err
	}

	exp, coin, err := cs.challengeRepo.GetChallengeRewards(challengeID)
	if err != nil {
		return err
	}

	err = cs.challengeRepo.UpdateRewardsGiven(challengeLogID)
	if err != nil {
		return err
	}

	err = cs.challengeRepo.AddUserRewards(userID, exp, coin)
	if err != nil {
		return err
	}

	return nil
}

func (cs *ChallengeService) GetActiveChallenges(userID string, page, perPage int, difficulty, title string) ([]challenges.ChallengeLog, int, error) {
	return cs.challengeRepo.GetChallengeLogByUserID(userID, page, perPage, difficulty, title)
}

func (cs *ChallengeService) GetUnclaimedChallenges(userID string, isAdmin bool, page, limit int, difficulty, title string) ([]challenges.Challenge, int, error) {
	return cs.challengeRepo.GetUnclaimedChallenges(userID, isAdmin, page, limit, difficulty, title)
}

func (cs *ChallengeService) GetChallengeDetailsWithConfirmations(userID, challengeLogID string) (challenges.ChallengeLogDetails, error) {
	challengeLog, err := cs.challengeRepo.GetChallengeLogByID(challengeLogID)
	if err != nil {
		return challenges.ChallengeLogDetails{}, err
	}

	if challengeLog.UserID != userID {
		return challenges.ChallengeLogDetails{}, constant.ErrUnauthorized
	}

	confirmations, err := cs.challengeRepo.GetConfirmationsByChallengeID(challengeLog.ChallengeID, userID)
	if err != nil {
		return challenges.ChallengeLogDetails{}, err
	}

	result := challenges.ChallengeLogDetails{
		ChallengeLog:  challengeLog,
		Confirmations: confirmations,
	}

	return result, nil
}

func (cs *ChallengeService) GetChallengeDetails(challengeID string) (challenges.ChallengeDetails, error) {
	challenge, err := cs.challengeRepo.GetChallengeByID(challengeID)
	if err != nil {
		return challenges.ChallengeDetails{}, err
	}

	tasks, err := cs.challengeRepo.GetTasksByChallengeIDforUser(challengeID)
	if err != nil {
		return challenges.ChallengeDetails{}, err
	}

	details := challenges.ChallengeDetails{
		ID:               challenge.ID,
		Title:            challenge.Title,
		Difficulty:       challenge.Difficulty,
		ActionCount:      challenge.ActionCount,
		ParticipantCount: challenge.ParticipantCount,
		ChallengeImg:     challenge.ChallengeImg,
		Description:      challenge.Description,
		DurationDays:     challenge.DurationDays,
		Exp:              challenge.Exp,
		Coin:             challenge.Coin,
		Tasks:            tasks,
	}

	return details, nil
}

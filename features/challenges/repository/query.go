package repository

import (
	"errors"
	"greenenvironment/constant"
	"greenenvironment/features/challenges"
	userRepo "greenenvironment/features/users/repository"
	"log"
	"time"

	"gorm.io/gorm"
)

type ChallengeData struct {
	DB *gorm.DB
}

func NewChallengeRepository(db *gorm.DB) challenges.ChallengeRepoInterface {
	return &ChallengeData{
		DB: db,
	}
}

func (cd *ChallengeData) Create(challenge challenges.Challenge) error {
	newChallenge := Challenge{
		ID:           challenge.ID,
		Author:       challenge.Author,
		Title:        challenge.Title,
		Difficulty:   challenge.Difficulty,
		ChallengeImg: challenge.ChallengeImg,
		Description:  challenge.Description,
		DurationDays: challenge.DurationDays,
		Exp:          challenge.Exp,
		Coin:         challenge.Coin,
	}

	for _, impactcategory := range challenge.ImpactCategories {
		newChallenge.ImpactCategories = append(newChallenge.ImpactCategories, ChallengeImpactCategory{
			ID:               impactcategory.ID,
			ChallengeID:      newChallenge.ID,
			ImpactCategoryID: impactcategory.ImpactCategoryID,
		})
	}
	err := cd.DB.Create(&newChallenge).Error
	if err != nil {
		return constant.ErrCreateChallenge
	}
	return nil
}

func (cd *ChallengeData) GetAllByPage(page int) ([]challenges.Challenge, int, error) {
	var challengeData []challenges.Challenge
	var totalChallenges int64

	err := cd.DB.Unscoped().Model(&Challenge{}).Count(&totalChallenges).Error
	if err != nil {
		return nil, 0, constant.ErrChallengeNotFound
	}

	challengesPerPage := 20
	totalPages := int((totalChallenges + int64(challengesPerPage) - 1) / int64(challengesPerPage))

	response := cd.DB.Unscoped().Preload("ImpactCategories.ImpactCategory").
		Select("challenges.*, challenges.author as author").
		Offset((page - 1) * challengesPerPage).
		Limit(challengesPerPage).
		Find(&challengeData)

	if response.Error != nil {
		return nil, 0, constant.ErrGetChallenge
	}

	if response.RowsAffected == 0 {
		return nil, 0, constant.ErrChallengeNotFound
	}

	return challengeData, totalPages, nil
}

func (cd *ChallengeData) GetByID(id string) (challenges.Challenge, error) {
	var challenge challenges.Challenge

	err := cd.DB.Model(&Challenge{}).
		Preload("ImpactCategories.ImpactCategory").
		Where("id = ?", id).
		Take(&challenge).Error

	if err != nil {
		return challenges.Challenge{}, constant.ErrChallengeNotFound
	}

	return challenge, nil
}

func (cd *ChallengeData) Update(challengeData challenges.Challenge) error {
	updatedChallenge := Challenge{
		ID:           challengeData.ID,
		Author:       challengeData.Author,
		Title:        challengeData.Title,
		Difficulty:   challengeData.Difficulty,
		ChallengeImg: challengeData.ChallengeImg,
		Description:  challengeData.Description,
		DurationDays: challengeData.DurationDays,
		Exp:          challengeData.Exp,
		Coin:         challengeData.Coin,
	}

	for _, impactCategory := range challengeData.ImpactCategories {
		updatedChallenge.ImpactCategories = append(updatedChallenge.ImpactCategories, ChallengeImpactCategory{
			ID:               impactCategory.ID,
			ChallengeID:      updatedChallenge.ID,
			ImpactCategoryID: impactCategory.ImpactCategoryID,
		})
	}

	tx := cd.DB.Begin()

	if err := tx.Where("challenge_id = ?", updatedChallenge.ID).Delete(&ChallengeImpactCategory{}).Error; err != nil {
		tx.Rollback()
		return constant.ErrUpdateChallenge
	}

	if err := tx.Model(&updatedChallenge).Where("id = ?", updatedChallenge.ID).Save(&updatedChallenge).Error; err != nil {
		tx.Rollback()
		return constant.ErrUpdateChallenge
	}

	return tx.Commit().Error
}

func (cd *ChallengeData) Delete(id string) error {
	tx := cd.DB.Begin()

	if err := tx.Where("challenge_id = ?", id).Delete(&ChallengeImpactCategory{}).Error; err != nil {
		tx.Rollback()
		return constant.ErrDeleteChallenge
	}

	if err := tx.Where("id = ?", id).Delete(&Challenge{}).Error; err != nil {
		tx.Rollback()
		return constant.ErrDeleteChallenge
	}

	return tx.Commit().Error
}

func (cd *ChallengeData) CreateTask(task challenges.ChallengeTask) error {
	newTask := ChallengeTask{
		ID:              task.ID,
		ChallengeID:     task.ChallengeID,
		DayNumber:       task.DayNumber,
		TaskDescription: task.TaskDescription,
	}
	err := cd.DB.Create(&newTask).Error
	if err != nil {
		return constant.ErrCreateTask
	}
	return nil
}

func (cd *ChallengeData) GetTasksByChallengeID(challengeID string) ([]challenges.ChallengeTask, error) {
	var tasks []ChallengeTask
	err := cd.DB.Where("challenge_id = ?", challengeID).Find(&tasks).Error
	if err != nil {
		return nil, constant.ErrFetchTasks
	}

	var result []challenges.ChallengeTask
	for _, task := range tasks {
		result = append(result, challenges.ChallengeTask{
			ID:              task.ID,
			ChallengeID:     task.ChallengeID,
			DayNumber:       task.DayNumber,
			TaskDescription: task.TaskDescription,
			CreatedAt:       task.CreatedAt,
			UpdatedAt:       task.UpdatedAt,
		})
	}
	return result, nil
}

func (cd *ChallengeData) GetTaskByID(taskID string) (challenges.ChallengeTask, error) {
	var task ChallengeTask
	err := cd.DB.Where("id = ?", taskID).First(&task).Error
	if err != nil {
		return challenges.ChallengeTask{}, constant.ErrTaskNotFound
	}
	return challenges.ChallengeTask{
		ID:              task.ID,
		ChallengeID:     task.ChallengeID,
		DayNumber:       task.DayNumber,
		TaskDescription: task.TaskDescription,
		CreatedAt:       task.CreatedAt,
		UpdatedAt:       task.UpdatedAt,
	}, nil
}

func (cd *ChallengeData) UpdateTask(task challenges.ChallengeTask) error {
	err := cd.DB.Model(&ChallengeTask{}).
		Where("id = ?", task.ID).
		Update("task_description", task.TaskDescription).
		Error
	if err != nil {
		return constant.ErrUpdateTask
	}
	return nil
}

func (cd *ChallengeData) DeleteTask(taskID string) error {
	err := cd.DB.Where("id = ?", taskID).Delete(&ChallengeTask{}).Error
	if err != nil {
		return constant.ErrDeleteTask
	}
	return nil
}

// User
func (cd *ChallengeData) CreateChallengeLog(log challenges.ChallengeLog) error {
	newLog := ChallengeLog{
		ID:           log.ID,
		ChallengeID:  log.ChallengeID,
		UserID:       log.UserID,
		Status:       log.Status,
		StartDate:    log.StartDate,
		Feed:         log.Feed,
		RewardsGiven: log.RewardsGiven,
	}

	err := cd.DB.Create(&newLog).Error
	if err != nil {
		return constant.ErrCreateChallengeLog
	}

	return nil
}

func (cd *ChallengeData) CreateChallengeConfirmation(confirmation challenges.ChallengeConfirmation) error {
	newConfirmation := ChallengeConfirmation{
		ID:              confirmation.ID,
		ChallengeTaskID: confirmation.ChallengeTaskID,
		UserID:          confirmation.UserID,
		Status:          confirmation.Status,
		ChallengeImg:    confirmation.ChallengeImg,
		SubmissionDate:  confirmation.SubmissionDate,
	}

	err := cd.DB.Create(&newConfirmation).Error
	if err != nil {
		return constant.ErrCreateChallengeConfirmation
	}

	return nil
}

func (cd *ChallengeData) IsChallengeTaken(userID, challengeID string) (bool, error) {
	var count int64
	err := cd.DB.Model(&ChallengeLog{}).
		Where("user_id = ? AND challenge_id = ?", userID, challengeID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (cd *ChallengeData) GetChallengeConfirmationByID(confirmationID string) (challenges.ChallengeConfirmation, error) {
	var confirmation ChallengeConfirmation
	err := cd.DB.Where("id = ?", confirmationID).First(&confirmation).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return challenges.ChallengeConfirmation{}, constant.ErrChallengeConfirmationNotFound
	}
	if err != nil {
		return challenges.ChallengeConfirmation{}, err
	}

	return challenges.ChallengeConfirmation{
		ID:              confirmation.ID,
		ChallengeTaskID: confirmation.ChallengeTaskID,
		UserID:          confirmation.UserID,
		Status:          confirmation.Status,
		ChallengeImg:    confirmation.ChallengeImg,
		SubmissionDate:  confirmation.SubmissionDate,
	}, nil
}

func (cd *ChallengeData) UpdateChallengeConfirmation(confirmation challenges.ChallengeConfirmation) error {
	err := cd.DB.Model(&ChallengeConfirmation{}).
		Where("id = ?", confirmation.ID).
		Updates(ChallengeConfirmation{
			Status:         confirmation.Status,
			ChallengeImg:   confirmation.ChallengeImg,
			SubmissionDate: confirmation.SubmissionDate,
		}).Error

	if err != nil {
		return constant.ErrUpdateChallengeConfirmation
	}

	return nil
}

func (cd *ChallengeData) GetChallengeTaskByID(taskID string) (challenges.ChallengeTask, error) {
	var task ChallengeTask
	err := cd.DB.Where("id = ?", taskID).First(&task).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return challenges.ChallengeTask{}, constant.ErrChallengeTaskNotFound
	}
	if err != nil {
		return challenges.ChallengeTask{}, err
	}

	return challenges.ChallengeTask{
		ID:              task.ID,
		ChallengeID:     task.ChallengeID,
		DayNumber:       task.DayNumber,
		TaskDescription: task.TaskDescription,
	}, nil
}

func (cd *ChallengeData) GetChallengeLogByChallengeIDAndUserID(challengeID, userID string) (challenges.ChallengeLog, error) {
	var log ChallengeLog
	err := cd.DB.Where("challenge_id = ? AND user_id = ?", challengeID, userID).First(&log).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return challenges.ChallengeLog{}, constant.ErrChallengeLogNotFound
	}
	if err != nil {
		return challenges.ChallengeLog{}, err
	}

	return challenges.ChallengeLog{
		ID:          log.ID,
		ChallengeID: log.ChallengeID,
		UserID:      log.UserID,
		Status:      log.Status,
		StartDate:   log.StartDate,
		Feed:        log.Feed,
	}, nil
}

func (cd *ChallengeData) GetConfirmationsByChallengeID(challengeID, userID string) ([]challenges.ChallengeConfirmation, error) {
	var confirmations []ChallengeConfirmation

	err := cd.DB.Table("challenge_confirmations").
		Select("challenge_confirmations.*, challenge_tasks.challenge_id").
		Joins("JOIN challenge_tasks ON challenge_tasks.id = challenge_confirmations.challenge_task_id").
		Where("challenge_tasks.challenge_id = ? AND challenge_confirmations.user_id = ?", challengeID, userID).
		Find(&confirmations).Error

	if err != nil {
		return nil, err
	}

	var result []challenges.ChallengeConfirmation
	for _, confirmation := range confirmations {
		result = append(result, challenges.ChallengeConfirmation{
			ID:              confirmation.ID,
			ChallengeTaskID: confirmation.ChallengeTaskID,
			UserID:          confirmation.UserID,
			Status:          confirmation.Status,
			ChallengeImg:    confirmation.ChallengeImg,
			SubmissionDate:  confirmation.SubmissionDate,
		})
	}

	return result, nil
}

func (cd *ChallengeData) UpdateChallengeLog(log challenges.ChallengeLog) error {
	err := cd.DB.Model(&ChallengeLog{}).
		Where("id = ?", log.ID).
		Updates(ChallengeLog{
			Status: log.Status,
		}).Error

	if err != nil {
		return constant.ErrUpdateChallengeLog
	}

	return nil
}

func (cd *ChallengeData) UpdateTaskAndChallengeStatus() error {
	now := time.Now()

	err := cd.DB.Model(&ChallengeConfirmation{}).
		Where("DATE_ADD(start_date, INTERVAL (day_number - 1) DAY) < ? AND status = ?", now, "Progress").
		Updates(map[string]interface{}{
			"status": "Failed",
		}).Error
	if err != nil {
		log.Printf("Error updating task statuses: %v", err)
		return err
	}

	rows, err := cd.DB.Raw(`
			UPDATE challenge_logs cl
			SET cl.status = 'Failed'
			WHERE cl.id IN (
					SELECT cc.challenge_log_id
					FROM challenge_confirmations cc
					WHERE cc.status = 'Failed'
					GROUP BY cc.challenge_log_id
					HAVING COUNT(*) = (
							SELECT COUNT(*)
							FROM challenge_confirmations
							WHERE challenge_log_id = cc.challenge_log_id
					)
			)
	`).Rows()
	if err != nil {
		log.Printf("Error updating challenge log statuses: %v", err)
		return err
	}
	defer rows.Close()

	log.Printf("Challenge and task statuses updated successfully.")
	return nil
}

func (cd *ChallengeData) IsRewardClaimed(challengeLogID string) (bool, error) {
	var log ChallengeLog
	err := cd.DB.Select("rewards_given").Where("id = ?", challengeLogID).First(&log).Error
	if err != nil {
		return false, err
	}
	return log.RewardsGiven, nil
}

func (cd *ChallengeData) UpdateRewardsGiven(challengeLogID string) error {
	err := cd.DB.Model(&ChallengeLog{}).Where("id = ?", challengeLogID).Update("rewards_given", true).Error
	return err
}

func (cd *ChallengeData) AddUserRewards(userID string, exp int, coin int) error {
	var user userRepo.User
	err := cd.DB.Model(&user).Where("id = ?", userID).Updates(map[string]interface{}{
		"exp":  gorm.Expr("exp + ?", exp),
		"coin": gorm.Expr("coin + ?", coin),
	}).Error
	return err
}

func (cd *ChallengeData) GetChallengeIDByLogID(challengeLogID string) (string, error) {
	var log ChallengeLog
	err := cd.DB.Select("challenge_id").Where("id = ?", challengeLogID).First(&log).Error
	if err != nil {
		return "", err
	}
	return log.ChallengeID, nil
}

func (cd *ChallengeData) GetChallengeRewards(challengeID string) (int, int, error) {
	var challenge Challenge
	err := cd.DB.Select("exp", "coin").Where("id = ?", challengeID).First(&challenge).Error
	if err != nil {
		return 0, 0, err
	}
	return challenge.Exp, challenge.Coin, nil
}

func (cd *ChallengeData) GetActiveChallengeLogByUserID(userID string, page, perPage int) ([]challenges.ChallengeLog, int, error) {
	var logs []struct {
		ChallengeLog
		Title        string
		Difficulty   string
		ChallengeImg string
		Description  string
		DurationDays int
		Exp          int
		Coin         int
	}
	var totalLogs int64

	err := cd.DB.Table("challenge_logs").
		Select("challenge_logs.*, challenges.title, challenges.difficulty, challenges.challenge_img, challenges.description, challenges.duration_days, challenges.exp, challenges.coin").
		Joins("JOIN challenges ON challenge_logs.challenge_id = challenges.id").
		Where("challenge_logs.user_id = ? AND challenge_logs.status = ?", userID, "Progress").
		Count(&totalLogs).
		Offset((page - 1) * perPage).
		Limit(perPage).
		Find(&logs).Error
	if err != nil {
		return nil, 0, err
	}

	var result []challenges.ChallengeLog
	for _, log := range logs {
		result = append(result, challenges.ChallengeLog{
			ID:          log.ID,
			ChallengeID: log.ChallengeID,
			UserID:      log.UserID,
			Status:      log.Status,
			StartDate:   log.StartDate,
			Feed:        log.Feed,
			Challenge: challenges.Challenge{
				Title:        log.Title,
				Difficulty:   log.Difficulty,
				ChallengeImg: log.ChallengeImg,
				Description:  log.Description,
				DurationDays: log.DurationDays,
				Exp:          log.Exp,
				Coin:         log.Coin,
			},
		})
	}

	totalPages := int((totalLogs + int64(perPage) - 1) / int64(perPage))
	return result, totalPages, nil
}

func (cd *ChallengeData) GetUnclaimedChallenges(userID string, isAdmin bool, page int, limit int) ([]challenges.Challenge, int, error) {
	var claimedChallenges []string
	err := cd.DB.Model(&ChallengeLog{}).
		Where("user_id = ?", userID).
		Pluck("challenge_id", &claimedChallenges).Error
	if err != nil {
		return nil, 0, err
	}

	var totalRecords int64
	query := cd.DB.Table("challenges").
		Select(`
			challenges.*, 
			COUNT(DISTINCT challenge_confirmations.id) AS action_count, 
			COUNT(DISTINCT challenge_logs.user_id) AS participant_count
		`).
		Joins("LEFT JOIN challenge_tasks ON challenge_tasks.challenge_id = challenges.id").
		Joins("LEFT JOIN challenge_confirmations ON challenge_confirmations.challenge_task_id = challenge_tasks.id").
		Joins("LEFT JOIN challenge_logs ON challenge_logs.challenge_id = challenges.id").
		Group("challenges.id").
		Count(&totalRecords)

	if !isAdmin {
		query = query.Where("challenges.deleted_at IS NULL")
	}

	if len(claimedChallenges) > 0 {
		query = query.Where("challenges.id NOT IN ?", claimedChallenges)
	}

	offset := (page - 1) * limit
	query = query.Offset(offset).Limit(limit)

	var challengeEntities []struct {
		challenges.Challenge
		ActionCount      int `gorm:"column:action_count"`
		ParticipantCount int `gorm:"column:participant_count"`
	}

	err = query.Find(&challengeEntities).Error
	if err != nil {
		return nil, 0, err
	}

	var result []challenges.Challenge
	for _, entity := range challengeEntities {
		result = append(result, challenges.Challenge{
			ID:               entity.ID,
			Author:           entity.Author,
			Title:            entity.Title,
			Difficulty:       entity.Difficulty,
			ChallengeImg:     entity.ChallengeImg,
			Description:      entity.Description,
			DurationDays:     entity.DurationDays,
			Exp:              entity.Exp,
			Coin:             entity.Coin,
			ActionCount:      entity.ActionCount,
			ParticipantCount: entity.ParticipantCount,
		})
	}

	totalPages := int((totalRecords + int64(limit) - 1) / int64(limit))
	return result, totalPages, nil
}

func (cd *ChallengeData) GetChallengeLogByID(challengeLogID string) (challenges.ChallengeLog, error) {
	var log ChallengeLog
	err := cd.DB.Where("id = ?", challengeLogID).First(&log).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return challenges.ChallengeLog{}, constant.ErrChallengeLogNotFound
	}
	if err != nil {
		return challenges.ChallengeLog{}, err
	}

	return challenges.ChallengeLog{
		ID:           log.ID,
		ChallengeID:  log.ChallengeID,
		UserID:       log.UserID,
		RewardsGiven: log.RewardsGiven,
		Status:       log.Status,
		StartDate:    log.StartDate,
		Feed:         log.Feed,
	}, nil
}

func (cd *ChallengeData) GetChallengeByID(challengeID string) (challenges.Challenge, error) {
	var challenge Challenge
	err := cd.DB.Where("id = ?", challengeID).First(&challenge).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return challenges.Challenge{}, constant.ErrChallengeNotFound
	}
	if err != nil {
		return challenges.Challenge{}, err
	}

	return challenges.Challenge{
		ID:           challenge.ID,
		Title:        challenge.Title,
		Difficulty:   challenge.Difficulty,
		ChallengeImg: challenge.ChallengeImg,
		Description:  challenge.Description,
		DurationDays: challenge.DurationDays,
		Exp:          challenge.Exp,
		Coin:         challenge.Coin,
	}, nil
}

func (cd *ChallengeData) GetTasksByChallengeIDforUser(challengeID string) ([]challenges.ChallengeTask, error) {
	var tasks []ChallengeTask
	err := cd.DB.Where("challenge_id = ?", challengeID).Order("day_number").Find(&tasks).Error
	if err != nil {
		return nil, err
	}

	var result []challenges.ChallengeTask
	for _, task := range tasks {
		result = append(result, challenges.ChallengeTask{
			ID:              task.ID,
			ChallengeID:     task.ChallengeID,
			DayNumber:       task.DayNumber,
			TaskDescription: task.TaskDescription,
		})
	}

	return result, nil
}

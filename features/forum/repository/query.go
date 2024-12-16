package repository

import (
	"errors"
	"greenenvironment/constant"
	"greenenvironment/features/forum"

	"gorm.io/gorm"
)

type ForumRepository struct {
	DB *gorm.DB
}

func NewForumRepository(db *gorm.DB) forum.ForumRepositoryInterface {
	return &ForumRepository{
		DB: db,
	}
}

// Forum
func (u *ForumRepository) GetAllForum() ([]forum.Forum, error) {
	var forum []forum.Forum
	if err := u.DB.Preload("User").Find(&forum).Error; err != nil {
		return nil, err
	}
	return forum, nil
}

func (u *ForumRepository) GetAllByPage(page int) ([]forum.ForumGetAll, int, error) {
	var forumData []forum.ForumGetAll

	var total int64
	count := u.DB.Model(&forum.Forum{}).Where("deleted_at IS NULL").Count(&total)
	if count.Error != nil {
		return nil, 0, constant.ErrProductEmpty
	}

	dataforumPerPage := 20
	totalPages := int((total + int64(dataforumPerPage) - 1) / int64(dataforumPerPage))

	tx := u.DB.Model(&Forum{}).Preload("User").Select("forums.*, (SELECT COUNT(*) FROM message_forums WHERE message_forums.forum_id = forums.id) AS message_count").Order("GREATEST(COALESCE(last_message_at, '2000-01-01'), created_at) DESC").Offset((page - 1) * dataforumPerPage).Limit(dataforumPerPage).Find(&forumData)
	if tx.Error != nil {
		return nil, 0, tx.Error
	}
	if tx.RowsAffected == 0 {
		return nil, 0, tx.Error
	}
	return forumData, totalPages, nil
}

func (u *ForumRepository) GetForumByID(ID string) (forum.Forum, error) {
	var forumData forum.Forum
	if err := u.DB.Model(&Forum{}).Preload("User").Where("id = ?", ID).First(&forumData).Error; err != nil {
		return forum.Forum{}, err
	}

	if err := u.DB.Model(&Forum{}).Where("id = ?", ID).Update("view", gorm.Expr("view + ?", 1)).Error; err != nil {
		return forum.Forum{}, err
	}

	return forumData, nil
}

func (u *ForumRepository) PostForum(forum forum.Forum) error {
	if err := u.DB.Create(&forum).Error; err != nil {
		return err
	}
	return nil
}

func (u *ForumRepository) UpdateForum(forumData forum.EditForum) error {
	var existingForum forum.Forum
	if err := u.DB.Where("id = ?", forumData.ID).First(&existingForum).Error; err != nil {
		return err
	}

	err := u.DB.Model(&existingForum).Updates(forumData).Error
	if err != nil {
		return err
	}
	return nil
}
func (u *ForumRepository) DeleteForum(forumID string) error {
	res := u.DB.Begin()

	result := res.Where("forum_id = ?", forumID).Delete(&MessageForum{})
	if result.Error != nil {
		res.Rollback()
		return result.Error
	}
	if err := res.Where("id = ?", forumID).Delete(&Forum{}); err.Error != nil {
		res.Rollback()
		return err.Error
	} else if err.RowsAffected == 0 {
		res.Rollback()
		return res.Error
	}

	return res.Commit().Error
}

func (u *ForumRepository) GetForumByUserID(userID string, page int) ([]forum.Forum, int, error) {
	var forum []forum.Forum
	var total int64

	dataforumPerPage := 20
	countResult := u.DB.Model(&Forum{}).Where("user_id = ?", userID).Count(&total)
	if countResult.Error != nil {
		return nil, 0, countResult.Error
	}

	if total == 0 {
		return nil, 0, errors.New("not found forum")
	}

	totalPages := int((total + int64(dataforumPerPage) - 1) / int64(dataforumPerPage))
	tx := u.DB.Model(&Forum{}).Preload("User").Order("created_at DESC, last_message_at DESC").Offset((page-1)*dataforumPerPage).Limit(dataforumPerPage).Where("user_id = ?", userID).Find(&forum)
	if tx.Error != nil {
		return nil, 0, constant.ErrGetProduct
	}
	return forum, totalPages, nil
}

// Message
func (u *ForumRepository) PostMessageForum(messageForum forum.MessageForum) error {
	res := u.DB.Begin()
	if err := res.Create(&messageForum).Error; err != nil {
		return err
	}

	if err := res.Model(&Forum{}).Where("id = ? AND deleted_at IS NULL", messageForum.ForumID).Update("last_message_at", messageForum.CreatedAt).Error; err != nil {
		res.Rollback()
		return err
	}

	var forum forum.Forum
	if err := u.DB.Where("id = ? AND deleted_at IS NULL", messageForum.ForumID).First(&forum).Error; err != nil {
		return errors.New("not found forum")
	}
	return res.Commit().Error
}

func (u *ForumRepository) DeleteMessageForum(messageID string) error {
	res := u.DB.Begin()

	if err := res.Where("id = ?", messageID).Delete(&MessageForum{}); err.Error != nil {
		res.Rollback()
		return errors.New("not found forum")
	} else if err.RowsAffected == 0 {
		res.Rollback()
		return errors.New("not found forum")
	}

	return res.Commit().Error
}

func (r *ForumRepository) GetMessageForumByID(ID string) (forum.MessageForum, error) {
	var messageForum forum.MessageForum
	if err := r.DB.Model(&MessageForum{}).Preload("User").Where("id = ? AND deleted_at IS NULL", ID).First(&messageForum).Error; err != nil {
		return forum.MessageForum{}, err
	}

	if err := r.DB.Model(&Forum{}).Where("id = ?", ID).Update("view", gorm.Expr("view + ?", 1)).Error; err != nil {
		return forum.MessageForum{}, err
	}
	return messageForum, nil
}

func (r *ForumRepository) UpdateMessageForum(message forum.EditMessage) error {
	var messageForum forum.MessageForum
	if err := r.DB.Where("id = ? AND deleted_at IS NULL", message.ID).First(&messageForum).Error; err != nil {
		return err
	}
	var forum forum.Forum
	if err := r.DB.Where("id = ? AND deleted_at IS NULL", messageForum.ForumID).First(&forum).Error; err != nil {
		return errors.New("not found forum")
	}
	err := r.DB.Model(&messageForum).Updates(message).Error
	if err != nil {
		return err
	}

	return nil
}

func (u *ForumRepository) GetMessagesByForumID(forumID string) ([]forum.MessageForum, error) {
	var messages []forum.MessageForum
	if err := u.DB.Model(&MessageForum{}).Preload("User").Where("forum_id = ? AND deleted_at IS NULL", forumID).Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}

func (u *ForumRepository) GetMessagesByForumIDWithPagination(forumID string, page int, pageSize int) ([]forum.MessageForum, error) {
	var messages []forum.MessageForum
	offset := (page - 1) * pageSize
	if err := u.DB.Where("forum_id = ?", forumID).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}

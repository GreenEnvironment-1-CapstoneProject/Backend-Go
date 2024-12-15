package repository

import (
	users "greenenvironment/features/users/repository"
	"time"

	"gorm.io/gorm"
)

type Forum struct {
	*gorm.Model
	ID            string     `gorm:"primary_key;type:varchar(50);not null;column:id;"`
	Title         string     `gorm:"type:varchar(255);not null;column:title"`
	TopicImage    string     `gorm:"column:topic_image"`
	View          int        `gorm:"default:0;column:view"`
	UserID        string     `gorm:"type:varchar(50);not null;column:user_id"`
	User          users.User `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Description   string     `gorm:"type:text;not null;column:description"`
	LastMessageAt time.Time  `gorm:"column:last_message_at"`
	CreatedAt     time.Time  `gorm:"column:created_at"`
}

type MessageForum struct {
	*gorm.Model
	ID           string     `gorm:"primary_key;type:varchar(50);not null;column:id"`
	ForumID      string     `gorm:"type:varchar(50);not null;column:forum_id"`
	Forum        Forum      `gorm:"foreignKey:ForumID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	MessageImage string     `gorm:"column:message_image"`
	Message      string     `gorm:"type:text;not null;column:message"`
	UserID       string     `gorm:"type:varchar(50);not null;column:user_id"`
	User         users.User `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (Forum) TableName() string {
	return "forums"
}

func (MessageForum) TableName() string {
	return "message_forums"
}

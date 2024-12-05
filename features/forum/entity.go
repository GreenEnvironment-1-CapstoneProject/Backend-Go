package forum

import (
	"greenenvironment/features/users"
	"time"

	"github.com/labstack/echo/v4"
)

type Forum struct {
	ID          string
	Title       string
	Description string
	UserID      string
	View        int
	TopicImage  string
	User        users.User
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type MessageForum struct {
	ID           string
	UserID       string
	User         users.User
	ForumID      string
	Message      string
	MessageImage string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type EditForum struct {
	ID          string
	Title       string
	UserID      string
	Description string
	TopicImage  string
	UpdatedAt   time.Time
}

type EditMessage struct {
	ID           string
	UserID       string
	Message      string
	MessageImage string
	UpdatedAt    time.Time
}

type ForumControllerInterface interface {
	GetAllForum(c echo.Context) error
	GetForumByID(c echo.Context) error
	PostForum(c echo.Context) error
	UpdateForum(c echo.Context) error
	DeleteForum(c echo.Context) error

	PostMessageForum(c echo.Context) error
	DeleteMessageForum(c echo.Context) error
	UpdateMessageForum(c echo.Context) error
	GetMessageForumByID(c echo.Context) error
	GetForumByUserID(c echo.Context) error
}

type ForumServiceInterface interface {
	GetAllForum() ([]Forum, error)
	GetAllByPage(page int) ([]Forum, int, error)
	GetForumByID(ID string) (Forum, error)
	PostForum(Forum) error
	GetForumByUserID(ID string, page int) ([]Forum, int, error)
	UpdateForum(EditForum) error
	DeleteForum(forumID string) error

	PostMessageForum(MessageForum) error
	DeleteMessageForum(messageID string) error
	GetMessageForumByID(ID string) (MessageForum, error)
	UpdateMessageForum(EditMessage) error
	GetMessagesByForumID(forumID string) ([]MessageForum, error)
	GetMessagesByForumIDWithPagination(forumID string, page int, pageSize int) ([]MessageForum, error)
}

type ForumRepositoryInterface interface {
	GetAllForum() ([]Forum, error)
	GetAllByPage(page int) ([]Forum, int, error)
	GetForumByID(ID string) (Forum, error)
	PostForum(Forum) error
	GetForumByUserID(ID string, page int) ([]Forum, int, error)

	UpdateForum(EditForum) error
	DeleteForum(forumID string) error

	PostMessageForum(MessageForum) error
	DeleteMessageForum(messageID string) error
	GetMessageForumByID(ID string) (MessageForum, error)
	UpdateMessageForum(EditMessage) error
	GetMessagesByForumID(forumID string) ([]MessageForum, error)
	GetMessagesByForumIDWithPagination(forumID string, page int, pageSize int) ([]MessageForum, error)
}

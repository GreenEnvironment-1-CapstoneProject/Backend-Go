package controller

type CreateForumRequest struct {
	Title       string `json:"title" form:"title" validate:"required"`
	Description string `json:"description" form:"description" validate:"required"`
}

type EditForumRequest struct {
	Title       string `json:"title" form:"title" validate:"required"`
	Description string `json:"description" form:"description"  validate:"required"`
}

type CreateMessageForumRequest struct {
	ForumID  string `json:"forum_id" form:"forum_id" validate:"required"`
	Messages string `json:"messages" form:"messages" validate:"required"`
}

type EditMessageRequest struct {
	ForumID  string `json:"forum_id" form:"forum_id" validate:"required"`
	Messages string `json:"messages" form:"messages" validate:"required"`
}

package controller

type ForumGetAllResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	View        int    `json:"views"`
	TopicImage  string `json:"topic_image"`
	Author      Author `json:"author"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type ForumGetDetailResponse struct {
	ID            string            `json:"id"`
	Title         string            `json:"title"`
	Description   string            `json:"description"`
	TopicImage    string            `json:"topic_image"`
	View          int               `json:"views"`
	CreatedAt     string            `json:"created_at"`
	UpdatedAt     string            `json:"updated_at"`
	Author        Author            `json:"author"`
	ForumMessages []MessageResponse `json:"forum_messages"`
}

type Author struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

type MessageResponse struct {
	ID           string        `json:"id"`
	User         AuthorMessage `json:"user"`
	Message      string        `json:"message"`
	MessageImage string        `json:"message_image"`
	CreatedAt    string        `json:"created_at"`
	UpdatedAt    string        `json:"updated_at"`
}

type AuthorMessage struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

type PaginatedResponse struct {
	Data  interface{} `json:"data"`
	Page  int         `json:"page"`
	Size  int         `json:"size"`
	Total int         `json:"total"`
}

type MetadataResponse struct {
	CurrentPage int `json:"current_page"`
	TotalPage   int `json:"total_page"`
}

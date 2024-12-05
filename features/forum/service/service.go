package service

import (
	"greenenvironment/constant"
	"greenenvironment/features/forum"
	"greenenvironment/helper"
	"strings"
)

type ForumService struct {
	forumRepo forum.ForumRepositoryInterface
}

func NewForumService(data forum.ForumRepositoryInterface) forum.ForumServiceInterface {
	return &ForumService{
		forumRepo: data,
	}
}

func (s *ForumService) GetAllForum() ([]forum.Forum, error) {
	return s.forumRepo.GetAllForum()
}

func (s *ForumService) GetAllByPage(page int) ([]forum.Forum, int, error) {
	return s.forumRepo.GetAllByPage(page)
}

func (s *ForumService) GetForumByID(ID string) (forum.Forum, error) {
	if ID == "" {
		return forum.Forum{}, constant.ErrFieldData
	}
	return s.forumRepo.GetForumByID(ID)
}

func (s *ForumService) PostForum(forum forum.Forum) error {
	if forum.Title == "" || forum.Description == "" {
		return constant.ErrFieldData
	}
	if !helper.IsValidInput(forum.Title) || !helper.IsValidInput(forum.Description) {
		return constant.ErrFieldData
	}
	return s.forumRepo.PostForum(forum)
}

func (s *ForumService) UpdateForum(forum forum.EditForum) error {
	if forum.ID == "" {
		return constant.ErrFieldData
	}
	if strings.TrimSpace(forum.Title) == "" || strings.TrimSpace(forum.Description) == "" {
		return constant.ErrFieldData
	}

	if !helper.IsValidInput(forum.Title) || !helper.IsValidInput(forum.Description) {
		return constant.ErrFieldData
	}
	return s.forumRepo.UpdateForum(forum)
}

func (s *ForumService) DeleteForum(forumID string) error {
	return s.forumRepo.DeleteForum(forumID)
}

func (s *ForumService) PostMessageForum(messageForum forum.MessageForum) error {
	if messageForum.ForumID == "" {
		return constant.ErrFieldData
	}
	if messageForum.Message == "" {
		return constant.ErrFieldData
	}

	if !helper.IsValidInput(messageForum.Message) {
		return constant.ErrFieldData
	}
	return s.forumRepo.PostMessageForum(messageForum)
}

func (s *ForumService) DeleteMessageForum(productId string) error {
	return s.forumRepo.DeleteMessageForum(productId)
}

func (s *ForumService) GetMessageForumByID(ID string) (forum.MessageForum, error) {
	if ID == "" {
		return forum.MessageForum{}, constant.ErrFieldData
	}
	return s.forumRepo.GetMessageForumByID(ID)
}

func (s *ForumService) UpdateMessageForum(message forum.EditMessage) error {
	if message.ID == "" {
		return constant.ErrFieldData
	}
	if message.Message == "" {
		return constant.ErrFieldData
	}

	if !helper.IsValidInput(message.Message) {
		return constant.ErrFieldData
	}
	return s.forumRepo.UpdateMessageForum(message)
}

func (s *ForumService) GetForumByUserID(ID string, page int) ([]forum.Forum, int, error) {
	return s.forumRepo.GetForumByUserID(ID, page)
}

func (s *ForumService) GetMessagesByForumID(forumID string) ([]forum.MessageForum, error) {
	return s.forumRepo.GetMessagesByForumID(forumID)
}

func (s *ForumService) GetMessagesByForumIDWithPagination(forumID string, page int, pageSize int) ([]forum.MessageForum, error) {
	return s.forumRepo.GetMessagesByForumIDWithPagination(forumID, page, pageSize)
}

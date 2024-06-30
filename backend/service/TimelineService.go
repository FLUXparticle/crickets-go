package service

import (
	"crickets-go/repository"
	"time"
)

type TimelineService struct {
	postRepository *repository.PostRepository
}

func NewTimelineService(postRepository *repository.PostRepository) *TimelineService {
	return &TimelineService{
		postRepository: postRepository,
	}
}

func (s *TimelineService) Post(creator *repository.User, content string) {
	s.postRepository.Save(&repository.Post{
		Creator:   creator,
		Content:   content,
		CreatedAt: time.Now(),
	})
}

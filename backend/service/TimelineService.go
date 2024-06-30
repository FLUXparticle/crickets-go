package service

import (
	"crickets-go/repository"
	"fmt"
	"sync"
	"time"
)

type TimelineService struct {
	postRepository         *repository.PostRepository
	subscriptionRepository *repository.SubscriptionRepository

	channelsMap map[int][]chan *repository.Post
	mu          sync.RWMutex
}

func NewTimelineService(postRepository *repository.PostRepository, subscriptionRepository *repository.SubscriptionRepository) *TimelineService {
	return &TimelineService{
		postRepository:         postRepository,
		subscriptionRepository: subscriptionRepository,
		channelsMap:            make(map[int][]chan *repository.Post),
	}
}

func (s *TimelineService) TimelineUpdates(subscriber *repository.User) chan *repository.Post {
	subscriptions := s.subscriptionRepository.FindBySubscriberID(subscriber.ID)

	creators := make([]*repository.User, len(subscriptions))
	for i, sub := range subscriptions {
		creators[i] = sub.Creator
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	aggregator := make(chan *repository.Post)

	for _, creator := range creators {
		ch := make(chan *repository.Post)
		go func() {
			for post := range ch {
				aggregator <- post
			}
		}()
		s.channelsMap[creator.ID] = append(s.channelsMap[creator.ID], ch)
	}

	return aggregator
}

func (s *TimelineService) Post(creator *repository.User, content string) {
	post := &repository.Post{
		Creator:   creator,
		Content:   content,
		CreatedAt: time.Now(),
	}

	s.postRepository.Save(post)

	s.mu.RLock()
	defer s.mu.RUnlock()
	if channels, found := s.channelsMap[creator.ID]; found {
		for _, ch := range channels {
			fmt.Println("Sending:", post.Content)
			ch <- post
		}
	}
}

func (s *TimelineService) Search(query string) []*repository.Post {
	return s.postRepository.FindByContentContains(query)
}

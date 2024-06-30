package service

import (
	"crickets-go/repository"
	"fmt"
	"time"
)

type TimelineService struct {
	postRepository         *repository.PostRepository
	subscriptionRepository *repository.SubscriptionRepository
	pubSub                 *PubSub
}

func NewTimelineService(postRepository *repository.PostRepository, subscriptionRepository *repository.SubscriptionRepository, pubSub *PubSub) *TimelineService {
	return &TimelineService{
		postRepository:         postRepository,
		subscriptionRepository: subscriptionRepository,
		pubSub:                 pubSub,
	}
}

func (s *TimelineService) TimelineUpdates(subscriber *repository.User) chan *repository.Post {
	subscriptions := s.subscriptionRepository.FindBySubscriberID(subscriber.ID)

	creators := make([]*repository.User, len(subscriptions))
	for i, sub := range subscriptions {
		creators[i] = sub.Creator
	}

	aggregator := make(chan *repository.Post)

	for _, creator := range creators {
		ch := s.pubSub.Subscribe(userID(creator))
		go func() {
			for post := range ch {
				aggregator <- post
			}
		}()
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
	s.pubSub.Publish(userID(creator), post)
}

func (s *TimelineService) Search(query string) []*repository.Post {
	return s.postRepository.FindByContentContains(query)
}

func userID(user *repository.User) string {
	return fmt.Sprintf("creator%d", user.ID)
}

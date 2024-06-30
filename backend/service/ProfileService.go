package service

import (
	"crickets-go/repository"
	"errors"
)

type ProfileService struct {
	userRepository         *repository.UserRepository
	subscriptionRepository *repository.SubscriptionRepository
}

func NewProfileService(userRepository *repository.UserRepository, subscriptionRepository *repository.SubscriptionRepository) *ProfileService {
	return &ProfileService{
		userRepository:         userRepository,
		subscriptionRepository: subscriptionRepository,
	}
}

func (s *ProfileService) SubscriberCount(creatorID int) int {
	return len(s.subscriptionRepository.FindByCreatorID(creatorID))
}

func (s *ProfileService) Subscribe(subscriber *repository.User, server string, creatorName string) (string, error) {
	if server == "" {
		creator := s.userRepository.FindByUsername(creatorName)
		if creator == nil {
			return "", errors.New("user not found")
		}
		s.subscriptionRepository.Save(&repository.Subscription{
			Creator:    creator,
			Subscriber: subscriber,
		})
	} else {
		panic("Not implemented")
	}
	return "success", nil
}

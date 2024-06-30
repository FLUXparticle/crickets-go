package service

import "crickets-go/repository"

type ProfileService struct {
	subscriptionRepository *repository.SubscriptionRepository
}

func NewProfileService(subscriptionRepository *repository.SubscriptionRepository) *ProfileService {
	return &ProfileService{
		subscriptionRepository: subscriptionRepository,
	}
}

func (s *ProfileService) SubscriberCount(creatorID int) int {
	return len(s.subscriptionRepository.FindByCreatorID(creatorID))
}

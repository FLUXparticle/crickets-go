package repository

import "crickets-go/data"

type SubscriptionRepository struct {
	subscriptions []*data.Subscription
}

func NewSubscriptionRepository() *SubscriptionRepository {
	return &SubscriptionRepository{
		// Wäre nicht unbedingt nätig, weil nil in den meisten Fällen auch als Slice funktioniert.
		subscriptions: make([]*data.Subscription, 0),
	}
}

func (s *SubscriptionRepository) FindAll() []*data.Subscription {
	return s.subscriptions
}

func (r *SubscriptionRepository) FindByCreatorID(creatorID int32) []*data.Subscription {
	result := make([]*data.Subscription, 0)

	for _, sub := range r.subscriptions {
		if sub.Creator.ID == creatorID {
			result = append(result, sub)
		}
	}

	return result
}

func (r *SubscriptionRepository) FindBySubscriberServerAndSubscriberID(subscriberServer string, subscriberID int32) []*data.Subscription {
	result := make([]*data.Subscription, 0)

	for _, sub := range r.subscriptions {
		if sub.Subscriber.Server == subscriberServer && sub.Subscriber.ID == subscriberID {
			result = append(result, sub)
		}
	}

	return result
}

func (r *SubscriptionRepository) Save(subscription *data.Subscription) {
	r.subscriptions = append(r.subscriptions, subscription)
}

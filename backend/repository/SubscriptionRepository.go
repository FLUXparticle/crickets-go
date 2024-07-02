package repository

type SubscriptionRepository struct {
	subscriptions []*Subscription
}

func NewSubscriptionRepository() *SubscriptionRepository {
	return &SubscriptionRepository{
		// Wäre nicht unbedingt nätig, weil nil in den meisten Fällen auch als Slice funktioniert.
		subscriptions: make([]*Subscription, 0),
	}
}

func (s *SubscriptionRepository) FindAll() []*Subscription {
	return s.subscriptions
}

func (r *SubscriptionRepository) FindByCreatorID(creatorID int32) []*Subscription {
	result := make([]*Subscription, 0)

	for _, sub := range r.subscriptions {
		if sub.Creator.ID == creatorID {
			result = append(result, sub)
		}
	}

	return result
}

func (r *SubscriptionRepository) FindBySubscriberServerAndSubscriberID(subscriberServer string, subscriberID int32) []*Subscription {
	result := make([]*Subscription, 0)

	for _, sub := range r.subscriptions {
		if sub.Subscriber.Server == subscriberServer && sub.Subscriber.ID == subscriberID {
			result = append(result, sub)
		}
	}

	return result
}

func (r *SubscriptionRepository) Save(subscription *Subscription) {
	r.subscriptions = append(r.subscriptions, subscription)
}

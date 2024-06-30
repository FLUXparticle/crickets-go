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

func (r *SubscriptionRepository) FindByCreatorID(creatorID int) []*Subscription {
	result := make([]*Subscription, 0)

	for _, sub := range r.subscriptions {
		if sub.creator.ID == creatorID {
			result = append(result, sub)
		}
	}

	return result
}

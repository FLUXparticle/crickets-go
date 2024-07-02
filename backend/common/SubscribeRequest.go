package common

import "crickets-go/repository"

type SubscribeRequest struct {
	Subscriber  *repository.User `json:"subscriber" binding:"required"`
	CreatorName string           `json:"creatorName" binding:"required"`
}

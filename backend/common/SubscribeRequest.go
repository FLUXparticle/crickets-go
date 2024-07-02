package common

import (
	"crickets-go/data"
)

type SubscribeRequest struct {
	Subscriber  *data.User `json:"subscriber" binding:"required"`
	CreatorName string     `json:"creatorName" binding:"required"`
}

package common

import "crickets-go/repository"

type SubscribeResponse struct {
	Error string           `json:"error,omitempty"`
	User  *repository.User `json:"user,omitempty"`
}

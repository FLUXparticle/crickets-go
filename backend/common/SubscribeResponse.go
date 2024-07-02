package common

import "crickets-go/data"

type SubscribeResponse struct {
	Error string     `json:"error,omitempty"`
	User  *data.User `json:"user,omitempty"`
}

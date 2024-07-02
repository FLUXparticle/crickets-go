package data

import "time"

type Post struct {
	Creator   *User
	Content   string `json:"content" binding:"required"`
	CreatedAt time.Time
}

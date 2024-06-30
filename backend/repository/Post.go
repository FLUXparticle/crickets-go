package repository

import "time"

type Post struct {
	Creator   *User
	Content   string `json:"content"`
	CreatedAt time.Time
}

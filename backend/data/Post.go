package data

import "time"

type Post struct {
	ID        int64 `json:"id"`
	Creator   *User
	Content   string `json:"content" binding:"required"`
	CreatedAt time.Time
	Likes     int `json:"likes"`
}

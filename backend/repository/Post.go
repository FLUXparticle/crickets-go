package repository

import "time"

type Post struct {
	Creator   *User
	Content   string
	CreatedAt time.Time
}

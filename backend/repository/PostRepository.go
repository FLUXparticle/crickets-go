package repository

import (
	"crickets-go/data"
	"strings"
)

type PostRepository struct {
	posts  []*data.Post
	nextID int64
}

func NewPostRepository() *PostRepository {
	return &PostRepository{
		posts:  make([]*data.Post, 0),
		nextID: 1,
	}
}

func (r *PostRepository) FindByID(id int64) *data.Post {
	for _, post := range r.posts {
		if post.ID == id {
			return post
		}
	}
	return nil
}

func (r *PostRepository) Save(post *data.Post) {
	post.ID = r.nextID
	r.nextID++
	r.posts = append(r.posts, post)
}

func (r *PostRepository) FindByContentContains(query string) []*data.Post {
	result := make([]*data.Post, 0)

	for _, post := range r.posts {
		if strings.Contains(post.Content, query) {
			result = append(result, post)
		}
	}

	return result
}

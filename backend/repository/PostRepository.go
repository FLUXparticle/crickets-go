package repository

import (
	"crickets-go/data"
	"strings"
)

type PostRepository struct {
	posts []*data.Post
}

func NewPostRepository() *PostRepository {
	return &PostRepository{
		posts: make([]*data.Post, 0),
	}
}

func (r *PostRepository) Save(post *data.Post) {
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

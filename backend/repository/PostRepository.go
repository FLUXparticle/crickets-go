package repository

import "strings"

type PostRepository struct {
	posts []*Post
}

func NewPostRepository() *PostRepository {
	return &PostRepository{
		posts: make([]*Post, 0),
	}
}

func (r *PostRepository) Save(post *Post) {
	r.posts = append(r.posts, post)
}

func (r *PostRepository) FindByContentContains(query string) []*Post {
	result := make([]*Post, 0)

	for _, post := range r.posts {
		if strings.Contains(post.Content, query) {
			result = append(result, post)
		}
	}

	return result
}

package repository

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

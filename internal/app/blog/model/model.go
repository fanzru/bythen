package model

import "time"

// Post represents a blog post
type Post struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	AuthorID  int64     `json:"author_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewPost represents the request payload for creating a new blog post
type NewPost struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// UpdatePost represents the request payload for updating a blog post
type UpdatePost struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

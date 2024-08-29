package model

import "time"

// Comment represents a comment on a blog post
type Comment struct {
	ID         int64     `json:"id"`
	PostID     int64     `json:"post_id"`
	AuthorName string    `json:"author_name"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
}

// NewComment represents the request payload for adding a comment to a blog post
type NewComment struct {
	AuthorName string `json:"author_name"`
	Content    string `json:"content"`
}

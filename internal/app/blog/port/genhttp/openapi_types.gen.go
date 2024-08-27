// Package genhttp provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.3.0 DO NOT EDIT.
package genhttp

import (
	"time"
)

// NewPost defines model for NewPost.
type NewPost struct {
	Content string `json:"content"`
	Title   string `json:"title"`
}

// Post defines model for Post.
type Post struct {
	AuthorId  *int       `json:"author_id,omitempty"`
	Content   *string    `json:"content,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	Id        *int       `json:"id,omitempty"`
	Title     *string    `json:"title,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

// UpdatePost defines model for UpdatePost.
type UpdatePost struct {
	Content *string `json:"content,omitempty"`
	Title   *string `json:"title,omitempty"`
}

// CreatePostJSONRequestBody defines body for CreatePost for application/json ContentType.
type CreatePostJSONRequestBody = NewPost

// UpdatePostJSONRequestBody defines body for UpdatePost for application/json ContentType.
type UpdatePostJSONRequestBody = UpdatePost

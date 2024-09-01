package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/fanzru/bythen/internal/app/blog/model"
)

type PostRepositoryImpl interface {
	CreatePost(ctx context.Context, post *model.Post) (int64, error)
	GetPostByID(ctx context.Context, id int64) (*model.Post, error)
	UpdatePost(ctx context.Context, post *model.Post) error
	DeletePost(ctx context.Context, id int64) error
	ListPosts(ctx context.Context) ([]*model.Post, error)
}

type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) CreatePost(ctx context.Context, post *model.Post) (int64, error) {
	query := "INSERT INTO posts (title, content, author_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?)"
	result, err := r.db.ExecContext(ctx, query, post.Title, post.Content, post.AuthorID, post.CreatedAt, post.UpdatedAt)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *PostRepository) GetPostByID(ctx context.Context, id int64) (*model.Post, error) {
	query := "SELECT id, title, content, author_id, created_at, updated_at FROM posts WHERE id = ?"
	row := r.db.QueryRowContext(ctx, query, id)

	var post model.Post
	if err := row.Scan(&post.ID, &post.Title, &post.Content, &post.AuthorID, &post.CreatedAt, &post.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("post not found")
		}
		return nil, err
	}
	return &post, nil
}

func (r *PostRepository) UpdatePost(ctx context.Context, post *model.Post) error {
	query := "UPDATE posts SET title = ?, content = ?, updated_at = ? WHERE id = ?"
	_, err := r.db.ExecContext(ctx, query, post.Title, post.Content, post.UpdatedAt, post.ID)
	return err
}

func (r *PostRepository) DeletePost(ctx context.Context, id int64) error {
	query := "DELETE FROM posts WHERE id = ?"
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *PostRepository) ListPosts(ctx context.Context) ([]*model.Post, error) {
	query := "SELECT id, title, content, author_id, created_at, updated_at FROM posts"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*model.Post
	for rows.Next() {
		var post model.Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.AuthorID, &post.CreatedAt, &post.UpdatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}
	return posts, nil
}

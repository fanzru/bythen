package repo

import (
	"context"
	"database/sql"

	"github.com/fanzru/bythen/internal/app/comment/model"
)

type CommentRepositoryImpl interface {
	CreateComment(ctx context.Context, comment *model.Comment) (int64, error)
	ListComments(ctx context.Context, postID int64) ([]*model.Comment, error)
}

type CommentRepository struct {
	db *sql.DB
}

func NewCommentRepository(db *sql.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

func (r *CommentRepository) CreateComment(ctx context.Context, comment *model.Comment) (int64, error) {
	query := `INSERT INTO comments (post_id, author_id, content, created_at) VALUES (?, ?, ?, ?)`
	result, err := r.db.ExecContext(ctx, query, comment.PostID, comment.AuthorID, comment.Content, comment.CreatedAt)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *CommentRepository) ListComments(ctx context.Context, postID int64) ([]*model.Comment, error) {
	query := `SELECT id, post_id, author_id, content, created_at FROM comments WHERE post_id = ?`
	rows, err := r.db.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*model.Comment
	for rows.Next() {
		var comment model.Comment
		if err := rows.Scan(&comment.ID, &comment.PostID, &comment.AuthorID, &comment.Content, &comment.CreatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, &comment)
	}
	return comments, nil
}

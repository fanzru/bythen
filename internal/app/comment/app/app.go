package app

import (
	"context"
	"time"

	"github.com/fanzru/bythen/internal/app/comment/model"
	"github.com/fanzru/bythen/internal/app/comment/repo"
)

type CommentServiceImpl interface {
	AddComment(ctx context.Context, postID int64, userID int64, request *model.CreateCommentRequest) (*model.CommentResponse, error)
	ListComments(ctx context.Context, postID int64) ([]*model.CommentResponse, error)
}

type CommentService struct {
	repo repo.CommentRepositoryImpl
}

func NewCommentService(repo repo.CommentRepositoryImpl) *CommentService {
	return &CommentService{repo: repo}
}

func (s *CommentService) AddComment(ctx context.Context, postID int64, userID int64, request *model.CreateCommentRequest) (*model.CommentResponse, error) {
	comment := &model.Comment{
		PostID:    postID,
		AuthorID:  userID,
		Content:   request.Content,
		CreatedAt: time.Now(),
	}

	id, err := s.repo.CreateComment(ctx, comment)
	if err != nil {
		return nil, err
	}

	return &model.CommentResponse{
		ID:        id,
		PostID:    comment.PostID,
		AuthorID:  comment.AuthorID,
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt,
	}, nil
}

func (s *CommentService) ListComments(ctx context.Context, postID int64) ([]*model.CommentResponse, error) {
	comments, err := s.repo.ListComments(ctx, postID)
	if err != nil {
		return nil, err
	}

	var commentResponses []*model.CommentResponse
	for _, comment := range comments {
		commentResponses = append(commentResponses, &model.CommentResponse{
			ID:        comment.ID,
			PostID:    comment.PostID,
			AuthorID:  comment.AuthorID,
			Content:   comment.Content,
			CreatedAt: comment.CreatedAt,
		})
	}
	return commentResponses, nil
}

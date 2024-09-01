package app

import (
	"context"
	"time"

	"github.com/fanzru/bythen/internal/app/blog/model"
	"github.com/fanzru/bythen/internal/app/blog/repo"
)

type PostServiceImpl interface {
	CreatePost(ctx context.Context, request *model.CreatePostRequest, authorID int64) (*model.PostResponse, error)
	GetPostByID(ctx context.Context, id int64) (*model.PostResponse, error)
	UpdatePost(ctx context.Context, id int64, request *model.UpdatePostRequest) (*model.PostResponse, error)
	DeletePost(ctx context.Context, id int64) error
	ListPosts(ctx context.Context) ([]*model.PostResponse, error)
}

type PostService struct {
	repo repo.PostRepositoryImpl
}

func NewPostService(repo repo.PostRepositoryImpl) *PostService {
	return &PostService{repo: repo}
}

func (s *PostService) CreatePost(ctx context.Context, request *model.CreatePostRequest, authorID int64) (*model.PostResponse, error) {
	post := &model.Post{
		Title:     request.Title,
		Content:   request.Content,
		AuthorID:  authorID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	id, err := s.repo.CreatePost(ctx, post)
	if err != nil {
		return nil, err
	}

	return &model.PostResponse{
		ID:        id,
		Title:     post.Title,
		Content:   post.Content,
		AuthorID:  post.AuthorID,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	}, nil
}

func (s *PostService) GetPostByID(ctx context.Context, id int64) (*model.PostResponse, error) {
	post, err := s.repo.GetPostByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &model.PostResponse{
		ID:        post.ID,
		Title:     post.Title,
		Content:   post.Content,
		AuthorID:  post.AuthorID,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	}, nil
}

func (s *PostService) UpdatePost(ctx context.Context, id int64, request *model.UpdatePostRequest) (*model.PostResponse, error) {
	post, err := s.repo.GetPostByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if request.Title != "" {
		post.Title = request.Title
	}
	if request.Content != "" {
		post.Content = request.Content
	}
	post.UpdatedAt = time.Now()

	if err := s.repo.UpdatePost(ctx, post); err != nil {
		return nil, err
	}

	return &model.PostResponse{
		ID:        post.ID,
		Title:     post.Title,
		Content:   post.Content,
		AuthorID:  post.AuthorID,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	}, nil
}

func (s *PostService) DeletePost(ctx context.Context, id int64) error {
	return s.repo.DeletePost(ctx, id)
}

func (s *PostService) ListPosts(ctx context.Context) ([]*model.PostResponse, error) {
	posts, err := s.repo.ListPosts(ctx)
	if err != nil {
		return nil, err
	}

	var postResponses []*model.PostResponse
	for _, post := range posts {
		postResponses = append(postResponses, &model.PostResponse{
			ID:        post.ID,
			Title:     post.Title,
			Content:   post.Content,
			AuthorID:  post.AuthorID,
			CreatedAt: post.CreatedAt,
			UpdatedAt: post.UpdatedAt,
		})
	}
	return postResponses, nil
}

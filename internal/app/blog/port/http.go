package port

import (
	"encoding/json"
	"net/http"

	"github.com/fanzru/bythen/internal/app/blog/app"
	"github.com/fanzru/bythen/internal/app/blog/model"
	"github.com/fanzru/bythen/internal/app/blog/port/genhttp"
	"github.com/fanzru/bythen/internal/common/response"
)

type PostHandler struct {
	service app.PostServiceImpl
}

func NewPostHandler(service app.PostServiceImpl) *PostHandler {
	return &PostHandler{service: service}
}

// CreatePost handles the creation of a new blog post
func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req genhttp.NewPost
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteErrorResponse(w, "Invalid request payload", err, http.StatusBadRequest)
		return
	}
	// Retrieve the userID from context
	userID, ok := ctx.Value("userID").(int64)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	resp, err := h.service.CreatePost(ctx, &model.CreatePostRequest{
		Title:   req.Title,
		Content: req.Content,
	}, userID)
	if err != nil {
		response.WriteErrorResponse(w, "Failed to create post", err, http.StatusInternalServerError)
		return
	}

	response.WriteSuccessResponse(w, resp, http.StatusCreated)
}

// ListPosts handles fetching a list of blog posts
func (h *PostHandler) ListPosts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	posts, err := h.service.ListPosts(ctx)
	if err != nil {
		response.WriteErrorResponse(w, "Failed to fetch posts", err, http.StatusInternalServerError)
		return
	}

	response.WriteSuccessResponse(w, posts, http.StatusOK)
}

// GetPostById handles fetching a single post by its ID
func (h *PostHandler) GetPostById(w http.ResponseWriter, r *http.Request, id int64) {
	ctx := r.Context()

	post, err := h.service.GetPostByID(ctx, id)
	if err != nil {
		response.WriteErrorResponse(w, "Post not found", err, http.StatusNotFound)
		return
	}

	response.WriteSuccessResponse(w, post, http.StatusOK)
}

// UpdatePost handles updating a blog post
func (h *PostHandler) UpdatePost(w http.ResponseWriter, r *http.Request, id int64) {
	ctx := r.Context()

	var req genhttp.UpdatePost
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteErrorResponse(w, "Invalid request payload", err, http.StatusBadRequest)
		return
	}

	post, err := h.service.UpdatePost(ctx, id, &model.UpdatePostRequest{
		Title:   req.Title,
		Content: req.Content,
	})
	if err != nil {
		response.WriteErrorResponse(w, "Failed to update post", err, http.StatusInternalServerError)
		return
	}

	response.WriteSuccessResponse(w, post, http.StatusOK)
}

// DeletePost handles deleting a blog post
func (h *PostHandler) DeletePost(w http.ResponseWriter, r *http.Request, id int64) {
	ctx := r.Context()

	if err := h.service.DeletePost(ctx, id); err != nil {
		response.WriteErrorResponse(w, "Failed to delete post", err, http.StatusInternalServerError)
		return
	}

	response.WriteSuccessResponse(w, nil, http.StatusNoContent)
}

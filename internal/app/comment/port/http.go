package port

import (
	"encoding/json"
	"net/http"

	"github.com/fanzru/bythen/internal/app/comment/app"
	"github.com/fanzru/bythen/internal/app/comment/model"
	"github.com/fanzru/bythen/internal/app/comment/port/genhttp"
	"github.com/fanzru/bythen/internal/common/response"
)

type CommentHandler struct {
	service app.CommentServiceImpl
}

func NewCommentHandler(service app.CommentServiceImpl) *CommentHandler {
	return &CommentHandler{service: service}
}

func (h *CommentHandler) AddComment(w http.ResponseWriter, r *http.Request, postID int64) {
	ctx := r.Context()

	var req genhttp.NewComment
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteErrorResponse(w, "Invalid request payload", err, http.StatusBadRequest)
		return
	}

	userID := ctx.Value("userID").(int64)
	resp, err := h.service.AddComment(ctx, postID, userID, &model.CreateCommentRequest{
		Content: req.Content,
	})
	if err != nil {
		response.WriteErrorResponse(w, "Failed to add comment", err, http.StatusInternalServerError)
		return
	}

	response.WriteSuccessResponse(w, resp, http.StatusCreated)
}

func (h *CommentHandler) ListComments(w http.ResponseWriter, r *http.Request, postID int64) {
	ctx := r.Context()

	comments, err := h.service.ListComments(ctx, postID)
	if err != nil {
		response.WriteErrorResponse(w, "Failed to fetch comments", err, http.StatusInternalServerError)
		return
	}

	response.WriteSuccessResponse(w, comments, http.StatusOK)
}

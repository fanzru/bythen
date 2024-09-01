package port

import (
	"encoding/json"
	"net/http"

	"github.com/fanzru/bythen/internal/app/user/app"
	"github.com/fanzru/bythen/internal/app/user/model"
	"github.com/fanzru/bythen/internal/app/user/port/genhttp"
	"github.com/fanzru/bythen/internal/common/response"
)

type UserHandler struct {
	service app.UserServiceImpl
}

func NewUserHandler(service app.UserServiceImpl) *UserHandler {
	return &UserHandler{service: service}
}

// RegisterUser handles user registration
func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req genhttp.UserRegistration
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	resp, err := h.service.CreateUser(ctx, &model.CreateUserRequest{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		response.WriteErrorResponse(w, "User creation failed", err, http.StatusInternalServerError)
		return
	}

	response.WriteSuccessResponse(w, resp, http.StatusCreated)
}

// LoginUser handles user login and returns a JWT token
func (h *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req genhttp.UserLogin
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	resp, err := h.service.LoginUser(ctx, &model.UserLoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		response.WriteErrorResponse(w, err.Error(), err, http.StatusUnauthorized)
		return
	}

	response.WriteSuccessResponse(w, resp, http.StatusOK)
}

// package response

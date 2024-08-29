package port

import (
	"context"

	"github.com/fanzru/bythen/internal/app/user/app"
	"github.com/fanzru/bythen/internal/app/user/port/genhttp"

	"github.com/fanzru/bythen/internal/app/user/model"
)

type UserHandler struct {
	service app.UserServiceImpl
}

func NewUserHandler(service app.UserServiceImpl) *UserHandler {
	return &UserHandler{service: service}
}

// RegisterUser handles user registration
func (h *UserHandler) RegisterUser(ctx context.Context, request genhttp.RegisterUserRequestObject) (genhttp.RegisterUserResponseObject, error) {
	_, err := h.service.CreateUser(ctx, &model.CreateUserRequest{
		Name:     request.Body.Name,
		Email:    request.Body.Email,
		Password: request.Body.Password,
	})
	if err != nil {
		return &genhttp.RegisterUser500Response{}, err
	}

	return &genhttp.RegisterUser201Response{}, nil
}

// LoginUser handles user login and returns a JWT token
func (h *UserHandler) LoginUser(ctx context.Context, request genhttp.LoginUserRequestObject) (genhttp.LoginUserResponseObject, error) {
	response, err := h.service.LoginUser(ctx, &model.UserLoginRequest{
		Email:    request.Body.Email,
		Password: request.Body.Password,
	})
	if err != nil {
		return &genhttp.LoginUser401Response{}, nil
	}

	return genhttp.LoginUser200JSONResponse{Token: &response.Token}, nil
}

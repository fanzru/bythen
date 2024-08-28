package app

import (
	"time"

	"github.com/fanzru/bythen/internal/app/user/model"
	"github.com/fanzru/bythen/internal/app/user/repo"
)

type UserServiceImpl interface {
	CreateUser(request *model.CreateUserRequest) (*model.CreateUserResponse, error)
}

type UserService struct {
	repo repo.UserRepository
}

func NewUserService(repo repo.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(request *model.CreateUserRequest) (*model.CreateUserResponse, error) {
	user := &model.User{
		Name:      request.Name,
		Email:     request.Email,
		Password:  hashPassword(request.Password), // Implement a hashPassword function
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	id, err := s.repo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	return &model.CreateUserResponse{
		ID:    id,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

// Add more business logic methods as needed

func hashPassword(password string) string {
	// Implement your password hashing logic here
	return password // This is just a placeholder
}

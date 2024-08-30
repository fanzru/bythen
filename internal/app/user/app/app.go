package app

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/fanzru/bythen/internal/app/user/model"
	"github.com/fanzru/bythen/internal/app/user/repo"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceImpl interface {
	CreateUser(ctx context.Context, request *model.CreateUserRequest) (*model.CreateUserResponse, error)
	LoginUser(ctx context.Context, request *model.UserLoginRequest) (*model.AuthTokenResponse, error)
}

type UserService struct {
	repo      repo.UserRepositoryImpl
	secretKey string
}

func NewUserService(repo repo.UserRepositoryImpl, secretKey string) *UserService {
	return &UserService{repo: repo, secretKey: secretKey}
}

// CreateUser handles the registration of a new user
func (s *UserService) CreateUser(ctx context.Context, request *model.CreateUserRequest) (*model.CreateUserResponse, error) {
	// Hash the user's password
	hashedPassword, err := hashPassword(request.Password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Name:      request.Name,
		Email:     request.Email,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Create the user in the database
	id, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return &model.CreateUserResponse{
		ID:    id,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

// LoginUser handles user authentication and returns a JWT token if successful
func (s *UserService) LoginUser(ctx context.Context, request *model.UserLoginRequest) (*model.AuthTokenResponse, error) {
	// Find the user by email
	user, err := s.repo.GetUserByEmail(ctx, request.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Compare the provided password with the stored hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Generate a JWT token
	token, err := s.generateJWT(user)
	if err != nil {
		return nil, err
	}

	return &model.AuthTokenResponse{
		Token: token,
	}, nil
}

// hashPassword hashes the user's password using bcrypt
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// generateJWT generates a JWT token for the authenticated user
func (s *UserService) generateJWT(user *model.User) (string, error) {
	// Define JWT claims
	claims := &model.JWTClaims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   fmt.Sprintf("%v", user.ID),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		}}

	// Create a new token with the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with a secret key
	tokenString, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

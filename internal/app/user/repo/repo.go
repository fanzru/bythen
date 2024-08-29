package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/fanzru/bythen/internal/app/user/model"
)

type UserRepositoryImpl interface {
	CreateUser(ctx context.Context, user *model.User) (int64, error)
	GetUserByID(ctx context.Context, id int64) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser inserts a new user into the database
func (r *UserRepository) CreateUser(ctx context.Context, user *model.User) (int64, error) {
	stmt, err := r.db.PrepareContext(ctx, "INSERT INTO users (name, email, password, created_at, updated_at) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, user.Name, user.Email, user.Password, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// GetUserByID retrieves a user from the database by ID
func (r *UserRepository) GetUserByID(ctx context.Context, id int64) (*model.User, error) {
	stmt, err := r.db.PrepareContext(ctx, "SELECT id, name, email, password, created_at, updated_at FROM users WHERE id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var user model.User
	if err := stmt.QueryRowContext(ctx, id).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

// GetUserByEmail retrieves a user from the database by email
func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	stmt, err := r.db.PrepareContext(ctx, "SELECT id, name, email, password, created_at, updated_at FROM users WHERE email = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var user model.User
	if err := stmt.QueryRowContext(ctx, email).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

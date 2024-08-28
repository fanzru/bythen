package repo

import (
	"database/sql"
	"errors"

	"github.com/fanzru/bythen/internal/app/user/model"
)

type UserRepositoryImpl interface {
	CreateUser(user *model.User) (int64, error)
	GetUserByID(id int64) (*model.User, error)
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(user *model.User) (int64, error) {
	query := "INSERT INTO users (name, email, password, created_at, updated_at) VALUES (?, ?, ?, ?, ?)"
	result, err := r.db.Exec(query, user.Name, user.Email, user.Password, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *UserRepository) GetUserByID(id int64) (*model.User, error) {
	query := "SELECT id, name, email, created_at, updated_at FROM users WHERE id = ?"
	row := r.db.QueryRow(query, id)

	var user model.User
	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

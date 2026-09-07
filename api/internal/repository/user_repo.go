package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/thegenggo/equipment-loan/api/internal/model"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) (int64, error) {
	const query = `
		INSERT INTO users (email, password_hash, name, role)
		VALUES (:email, :password_hash, :name, :role)
	`

	result, err := r.db.NamedExecContext(ctx, query, user)
	if err != nil {
		return 0, fmt.Errorf("insert user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("last insert id: %w", err)
	}

	return id, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	const query = `
		SELECT id, email, password_hash, name, role
		FROM users
		WHERE email = ?`

	var user model.User
	if err := r.db.GetContext(ctx, &user, query, email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("select user by email: %w", err)
	}

	return &user, nil
}

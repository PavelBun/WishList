package repository

import (
	"context"
	"errors"
	"wishlist-api/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// UserPostgres implements UserRepository using PostgreSQL.
type UserPostgres struct {
	pool *pgxpool.Pool
}

// NewUserRepository creates a new UserPostgres repository.
func NewUserRepository(pool *pgxpool.Pool) *UserPostgres {
	return &UserPostgres{pool: pool}
}

// Create inserts a new user into the database.
func (r *UserPostgres) Create(ctx context.Context, email, passwordHash string) (*models.User, error) {
	var user models.User
	err := r.pool.QueryRow(ctx,
		"INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id, email, created_at, updated_at",
		email, passwordHash,
	).Scan(&user.ID, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail retrieves a user by email address.
func (r *UserPostgres) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}
	err := r.pool.QueryRow(ctx,
		"SELECT id, email, password_hash, created_at, updated_at FROM users WHERE email = $1",
		email,
	).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

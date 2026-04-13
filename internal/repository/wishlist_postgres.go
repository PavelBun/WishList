package repository

import (
	"context"
	"errors"
	"time"
	"wishlist-api/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// WishlistPostgres implements WishlistRepository using PostgreSQL.
type WishlistPostgres struct {
	pool *pgxpool.Pool
}

// NewWishlistRepository creates a new WishlistPostgres repository.
func NewWishlistRepository(pool *pgxpool.Pool) *WishlistPostgres {
	return &WishlistPostgres{pool: pool}
}

// Create inserts a new wishlist into the database.
func (r *WishlistPostgres) Create(ctx context.Context, userID int, title, description string, eventDate time.Time) (*models.Wishlist, error) {
	var w models.Wishlist
	err := r.pool.QueryRow(ctx,
		`INSERT INTO wishlists (user_id, title, description, event_date) 
         VALUES ($1, $2, $3, $4) 
         RETURNING id, user_id, title, description, event_date, access_token, created_at, updated_at`,
		userID, title, description, eventDate,
	).Scan(&w.ID, &w.UserID, &w.Title, &w.Description, &w.EventDate, &w.AccessToken, &w.CreatedAt, &w.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &w, nil
}

// GetByID retrieves a wishlist by its ID.
func (r *WishlistPostgres) GetByID(ctx context.Context, id int) (*models.Wishlist, error) {
	w := &models.Wishlist{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, user_id, title, description, event_date, access_token, created_at, updated_at 
         FROM wishlists WHERE id = $1`,
		id,
	).Scan(&w.ID, &w.UserID, &w.Title, &w.Description, &w.EventDate, &w.AccessToken, &w.CreatedAt, &w.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return w, nil
}

// GetByAccessToken retrieves a wishlist by its public access token.
func (r *WishlistPostgres) GetByAccessToken(ctx context.Context, token uuid.UUID) (*models.Wishlist, error) {
	w := &models.Wishlist{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, user_id, title, description, event_date, access_token, created_at, updated_at 
         FROM wishlists WHERE access_token = $1`,
		token,
	).Scan(&w.ID, &w.UserID, &w.Title, &w.Description, &w.EventDate, &w.AccessToken, &w.CreatedAt, &w.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return w, nil
}

// GetAllByUser returns all wishlists belonging to the given user.
func (r *WishlistPostgres) GetAllByUser(ctx context.Context, userID int) ([]models.Wishlist, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, title, description, event_date, access_token, created_at, updated_at 
         FROM wishlists WHERE user_id = $1 ORDER BY event_date DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var wishlists []models.Wishlist
	for rows.Next() {
		var w models.Wishlist
		err := rows.Scan(&w.ID, &w.UserID, &w.Title, &w.Description, &w.EventDate, &w.AccessToken, &w.CreatedAt, &w.UpdatedAt)
		if err != nil {
			return nil, err
		}
		wishlists = append(wishlists, w)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return wishlists, nil
}

// Update modifies an existing wishlist.
func (r *WishlistPostgres) Update(ctx context.Context, w *models.Wishlist) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE wishlists SET title = $1, description = $2, event_date = $3, updated_at = NOW() 
         WHERE id = $4`,
		w.Title, w.Description, w.EventDate, w.ID,
	)
	return err
}

// Delete removes a wishlist by ID.
func (r *WishlistPostgres) Delete(ctx context.Context, id int) error {
	_, err := r.pool.Exec(ctx, "DELETE FROM wishlists WHERE id = $1", id)
	return err
}

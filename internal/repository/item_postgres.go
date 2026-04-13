package repository

import (
	"context"
	"errors"
	"fmt"
	"wishlist-api/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ErrAlreadyBooked is returned when trying to book an already booked item.
var ErrAlreadyBooked = errors.New("item already booked")

// ItemPostgres implements ItemRepository using PostgreSQL.
type ItemPostgres struct {
	pool *pgxpool.Pool
}

// NewItemRepository creates a new ItemPostgres repository.
func NewItemRepository(pool *pgxpool.Pool) *ItemPostgres {
	return &ItemPostgres{pool: pool}
}

// Create inserts a new item into the database.
func (r *ItemPostgres) Create(ctx context.Context, item *models.Item) error {
	err := r.pool.QueryRow(ctx,
		`INSERT INTO items (wishlist_id, title, description, product_link, priority) 
         VALUES ($1, $2, $3, $4, $5) 
         RETURNING id, created_at, updated_at`,
		item.WishlistID, item.Title, item.Description, item.ProductLink, int(item.Priority),
	).Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		return fmt.Errorf("insert item: %w", err)
	}
	return nil
}

// GetByID retrieves an item by its UUID.
func (r *ItemPostgres) GetByID(ctx context.Context, id uuid.UUID) (*models.Item, error) {
	item := &models.Item{}
	var priorityInt int
	err := r.pool.QueryRow(ctx,
		`SELECT id, wishlist_id, title, description, product_link, priority, is_booked, created_at, updated_at 
         FROM items WHERE id = $1`,
		id,
	).Scan(&item.ID, &item.WishlistID, &item.Title, &item.Description, &item.ProductLink,
		&priorityInt, &item.IsBooked, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get item by id: %w", err)
	}
	item.Priority = models.Priority(priorityInt)
	return item, nil
}

// GetAllByWishlistID returns all items belonging to the given wishlist.
func (r *ItemPostgres) GetAllByWishlistID(ctx context.Context, wishlistID uuid.UUID) ([]models.Item, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, wishlist_id, title, description, product_link, priority, is_booked, created_at, updated_at 
         FROM items WHERE wishlist_id = $1 ORDER BY priority DESC, created_at ASC`,
		wishlistID,
	)
	if err != nil {
		return nil, fmt.Errorf("query items by wishlist: %w", err)
	}
	defer rows.Close()

	var items []models.Item
	for rows.Next() {
		var i models.Item
		var priorityInt int
		err := rows.Scan(&i.ID, &i.WishlistID, &i.Title, &i.Description, &i.ProductLink,
			&priorityInt, &i.IsBooked, &i.CreatedAt, &i.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("scan item row: %w", err)
		}
		i.Priority = models.Priority(priorityInt)
		items = append(items, i)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}
	return items, nil
}

// Update modifies an existing item.
func (r *ItemPostgres) Update(ctx context.Context, item *models.Item) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE items SET title = $1, description = $2, product_link = $3, priority = $4, updated_at = NOW() 
         WHERE id = $5`,
		item.Title, item.Description, item.ProductLink, int(item.Priority), item.ID,
	)
	if err != nil {
		return fmt.Errorf("update item: %w", err)
	}
	return nil
}

// Delete removes an item by ID.
func (r *ItemPostgres) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, "DELETE FROM items WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("delete item: %w", err)
	}
	return nil
}

// BookItem marks an item as booked. Returns ErrAlreadyBooked if already booked.
func (r *ItemPostgres) BookItem(ctx context.Context, id uuid.UUID) error {
	tag, err := r.pool.Exec(ctx,
		"UPDATE items SET is_booked = TRUE WHERE id = $1 AND is_booked = FALSE",
		id,
	)
	if err != nil {
		return fmt.Errorf("book item exec: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrAlreadyBooked
	}
	return nil
}

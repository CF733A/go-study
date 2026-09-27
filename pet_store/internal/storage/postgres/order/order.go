package order

import (
	"context"
	"fmt"
	"pet-store/internal/entities"

	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepository struct {
	db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{
		db: db,
	}
}

func (or *OrderRepository) Create(ctx context.Context, order *entities.Order) error {
	query := `
		INSERT INTO orders (pet_id, user_id, quantity, ship_date, status, complete)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	err := or.db.QueryRow(ctx, query,
		order.PetID,
		order.UserID,
		order.Quantity,
		order.ShipDate,
		order.Status,
		order.Complete).Scan(&order.ID)

	if err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}

	return nil
}

func (or *OrderRepository) GetByID(ctx context.Context, id int64) (*entities.Order, error) {
	query := `
		SELECT id, pet_id, user_id, quantity, ship_date, status, complete
		FROM orders 
		WHERE id = $1
	`

	order := &entities.Order{}
	err := or.db.QueryRow(ctx, query, id).Scan(
		&order.ID,
		&order.PetID,
		&order.UserID,
		&order.Quantity,
		&order.ShipDate,
		&order.Status,
		&order.Complete,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get order by ID: %w", err)
	}

	return order, nil
}

func (or *OrderRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM orders WHERE id = $1`

	_, err := or.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete order: %w", err)
	}

	return nil
}

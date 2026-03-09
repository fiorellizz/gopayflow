package database

import (
	"context"
	"database/sql"

	"github.com/fiorellizz/gopayflow/internal/domain"
)

type PostgresOrderRepository struct {
	db *sql.DB
}

func NewPostgresOrderRepository(db *sql.DB) *PostgresOrderRepository {
	return &PostgresOrderRepository{
		db: db,
	}
}

func (r *PostgresOrderRepository) Create(ctx context.Context, order *domain.Order) error {
	query := `
		INSERT INTO orders (amount, status, created_at)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	return r.db.QueryRowContext(
		ctx,
		query,
		order.Amount,
		order.Status,
		order.CreatedAt,
	).Scan(&order.ID)
}

func (r *PostgresOrderRepository) FindByID(ctx context.Context, id string) (*domain.Order, error) {
	query := `
		SELECT id, amount, status, created_at
		FROM orders
		WHERE id = $1
	`

	var order domain.Order

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&order.ID,
		&order.Amount,
		&order.Status,
		&order.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (r *PostgresOrderRepository) FindAll(ctx context.Context) ([]*domain.Order, error) {
	query := `
		SELECT id, amount, status, created_at
		FROM orders
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*domain.Order

	for rows.Next() {
		var order domain.Order

		err := rows.Scan(
			&order.ID,
			&order.Amount,
			&order.Status,
			&order.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		orders = append(orders, &order)
	}

	return orders, nil
}

func (r *PostgresOrderRepository) UpdateStatus(
	ctx context.Context,
	id string,
	status domain.OrderStatus,
) error {

	query := `
		UPDATE orders
		SET status = $1
		WHERE id = $2
	`

	_, err := r.db.ExecContext(ctx, query, status, id)

	return err
}

package domain

import "context"

type OrderRepository interface {
	Create(ctx context.Context, order *Order) error
	FindByID(ctx context.Context, id string) (*Order, error)
	FindAll(ctx context.Context) ([]*Order, error)
	UpdateStatus(ctx context.Context, id string, status OrderStatus) error
}
package domain

import "context"

type EventPublisher interface {
	PublishOrderCreated(ctx context.Context, order *Order) error
}
package storage

import (
	"context"
	"errors"
	"order-service/pkg/api/test"
)

type Interface interface {
	Add(ctx context.Context, order *test.Order) error
	Get(ctx context.Context, id string) (*test.Order, error)
	Update(ctx context.Context, order *test.Order) (*test.Order, error)
	Delete(ctx context.Context, id string) error
	GetAll(ctx context.Context) ([]*test.Order, error)
	Stop()
}

var (
	ValueNotFound = errors.New("value not found")
)

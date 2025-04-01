package repository

import (
	"context"
	"fmt"
	"order-service/internal/storage"
	"order-service/internal/storage/postgresql"
	"order-service/pkg/api/test"
)

type Repository struct {
	storage storage.Interface
}

func New(storage *postgresql.Storage) *Repository {
	return &Repository{storage: storage}
}

func (r *Repository) Add(ctx context.Context, order *test.Order) error {
	if err := r.storage.Add(ctx, order); err != nil {
		return fmt.Errorf("failed to add order: %w", err)
	}
	return nil
}

func (r *Repository) Get(ctx context.Context, id string) (*test.Order, error) {
	order, err := r.storage.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}
	return order, nil
}

func (r *Repository) Update(ctx context.Context, order *test.Order) (*test.Order, error) {
	newOrder, err := r.storage.Update(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("failed to update order: %w", err)
	}
	return newOrder, nil
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	if err := r.storage.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete order: %w", err)
	}
	return nil
}

func (r *Repository) GetAll(ctx context.Context) ([]*test.Order, error) {
	orders, err := r.storage.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all orders: %w", err)
	}
	return orders, nil
}

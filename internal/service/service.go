package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"order-service/internal/logger"
	_ "order-service/internal/logger"
	"order-service/internal/storage"
	"order-service/pkg/api/test"
)

type OrderRepository interface {
	Add(ctx context.Context, order *test.Order) error
	Get(ctx context.Context, id string) (*test.Order, error)
	Update(ctx context.Context, order *test.Order) (*test.Order, error)
	Delete(ctx context.Context, id string) error
	GetAll(ctx context.Context) ([]*test.Order, error)
}

type Service struct {
	OrderRepository OrderRepository
	test.UnimplementedOrderServiceServer
}

func Register(gRPC *grpc.Server, repo OrderRepository) {
	service := &Service{
		OrderRepository: repo,
	}
	test.RegisterOrderServiceServer(gRPC, service)
}

func (s *Service) CreateOrder(ctx context.Context,
	req *test.CreateOrderRequest) (*test.CreateOrderResponse, error) {
	id := uuid.New().String()
	order := &test.Order{
		Id:       id,
		Item:     req.GetItem(),
		Quantity: req.GetQuantity(),
	}
	err := s.OrderRepository.Add(ctx, order)
	if err != nil {
		return nil, err
	}
	logger.GetLogger().Info(ctx, "order is add", zap.String("id", id))
	return &test.CreateOrderResponse{Id: id}, nil
}

func (s *Service) GetOrder(ctx context.Context,
	req *test.GetOrderRequest) (*test.GetOrderResponse, error) {

	order, err := s.OrderRepository.Get(ctx, req.GetId())
	if err != nil {
		if errors.Is(err, storage.ValueNotFound) {
			logger.GetLogger().Error(ctx, "order not found", zap.String("id", req.GetId()))
			return nil, status.Errorf(codes.NotFound, "order with ID %s not found", req.GetId())
		}
		logger.GetLogger().Error(ctx, "get order error", zap.String("id", req.GetId()), zap.Error(err))
		return nil, err
	}
	logger.GetLogger().Info(ctx, "order is get", zap.String("id", req.GetId()))
	return &test.GetOrderResponse{Order: order}, nil
}
func (s *Service) ListOrders(ctx context.Context,
	req *test.ListOrdersRequest) (*test.ListOrdersResponse, error) {

	orders, err := s.OrderRepository.GetAll(ctx)
	if err != nil {
		logger.GetLogger().Error(ctx, "get all orders error", zap.Error(err))
		return nil, err
	}
	logger.GetLogger().Info(ctx, "get all orders", zap.Int("count", len(orders)))
	return &test.ListOrdersResponse{Orders: orders}, nil
}

func (s *Service) UpdateOrder(ctx context.Context,
	req *test.UpdateOrderRequest) (*test.UpdateOrderResponse, error) {
	order := &test.Order{
		Id:       req.GetId(),
		Item:     req.GetItem(),
		Quantity: req.GetQuantity(),
	}
	newOrder, err := s.OrderRepository.Update(ctx, order)
	if err != nil {
		if errors.Is(err, storage.ValueNotFound) {
			logger.GetLogger().Error(ctx, "order not found", zap.String("id", req.GetId()))
			return nil, status.Errorf(codes.NotFound, "order with ID %s not found", req.GetId())
		}
		logger.GetLogger().Error(ctx, "update order error", zap.String("id", req.GetId()), zap.Error(err))
		return nil, err
	}
	logger.GetLogger().Info(ctx, "order is updated", zap.String("id", req.GetId()))
	return &test.UpdateOrderResponse{Order: newOrder}, nil
}

func (s *Service) DeleteOrder(ctx context.Context,
	req *test.DeleteOrderRequest) (*test.DeleteOrderResponse, error) {

	err := s.OrderRepository.Delete(ctx, req.GetId())
	if err != nil {
		logger.GetLogger().Info(ctx, "Order successfully deleted", zap.String("id", req.GetId()))
		return &test.DeleteOrderResponse{Success: true}, nil
	}
	logger.GetLogger().Error(ctx, "Failed to delete order", zap.String("id", req.GetId()))
	return &test.DeleteOrderResponse{Success: false}, err
}

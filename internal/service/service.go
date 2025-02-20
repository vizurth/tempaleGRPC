package orderservice

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"sync"
	test "tempaleGRCP/pkg/api/test/api"
	"tempaleGRCP/pkg/logger"
)

var (
	mu    sync.Mutex
	store = make(map[string]test.Order)
)

type Service struct {
	test.OrderServiceServer
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) CreateOrder(ctx context.Context, req *test.CreateOrderRequest) (*test.CreateOrderResponse, error) {
	uniqueID := uuid.New().String()
	mu.Lock()
	store[uniqueID] = test.Order{Id: uniqueID, Item: req.Item, Quantity: req.Quantity}
	mu.Unlock()
	logger.GetLoggerFromCtx(ctx).Info(ctx, "CreateOrder called")
	return &test.CreateOrderResponse{
		Id: uniqueID,
	}, nil
}

func (s *Service) GetOrder(ctx context.Context, req *test.GetOrderRequest) (*test.GetOrderResponse, error) {
	mu.Lock()
	order, ok := store[req.Id]
	mu.Unlock()
	if !ok {
		logger.GetLoggerFromCtx(ctx).Info(ctx, "order not found", zap.String("id", req.Id))
		return &test.GetOrderResponse{Order: nil}, fmt.Errorf("order not found")
	}
	logger.GetLoggerFromCtx(ctx).Info(ctx, "GetOrder called", zap.String("id", req.Id))
	return &test.GetOrderResponse{Order: &order}, nil
}

func (s *Service) UpdateOrder(ctx context.Context, req *test.UpdateOrderRequest) (*test.UpdateOrderResponse, error) {
	mu.Lock()
	order, ok := store[req.Id]
	defer mu.Unlock()
	if !ok {
		logger.GetLoggerFromCtx(ctx).Info(ctx, "order not found", zap.String("id", req.Id))
		return &test.UpdateOrderResponse{Order: nil}, fmt.Errorf("order not found")
	}
	logger.GetLoggerFromCtx(ctx).Info(ctx, "UpdateOrder called", zap.String("id", req.Id))
	order.Item = req.Item
	order.Quantity = req.Quantity
	store[req.Id] = order

	return &test.UpdateOrderResponse{Order: &order}, nil

}

func (s *Service) DeleteOrder(ctx context.Context, req *test.DeleteOrderRequest) (*test.DeleteOrderResponse, error) {
	mu.Lock()
	_, ok := store[req.Id]
	defer mu.Unlock()
	if !ok {
		logger.GetLoggerFromCtx(ctx).Info(ctx, "order not found", zap.String("id", req.Id))
		return &test.DeleteOrderResponse{Success: false}, fmt.Errorf("order not found")
	}
	delete(store, req.Id)
	logger.GetLoggerFromCtx(ctx).Info(ctx, "DeleteOrder called", zap.String("id", req.Id))
	return &test.DeleteOrderResponse{Success: true}, nil
}

func (s *Service) ListOrders(ctx context.Context, req *test.ListOrdersRequest) (*test.ListOrdersResponse, error) {
	mu.Lock()
	defer mu.Unlock()
	orders := make([]*test.Order, 0, len(store))
	for _, order := range store {
		orders = append(orders, &order)
	}
	logger.GetLoggerFromCtx(ctx).Info(ctx, "DeleteOrder called")
	return &test.ListOrdersResponse{Orders: orders}, nil
}

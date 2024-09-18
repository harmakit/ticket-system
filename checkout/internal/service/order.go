package service

import (
	"context"
	"ticket-system/checkout/internal/model"
	"ticket-system/checkout/internal/repository"
)

type orderService struct {
	orderRepository repository.OrderRepository
}

func NewOrderService(eventRepository repository.OrderRepository) OrderService {
	return &orderService{
		eventRepository,
	}
}

func (s orderService) GetOrder(ctx context.Context, id model.UUID) (*model.Order, error) {
	return s.orderRepository.Find(ctx, id)
}

func (s orderService) CreateOrder(ctx context.Context, o *model.Order) error {
	return s.orderRepository.Create(ctx, o)
}

func (s orderService) UpdateOrder(ctx context.Context, order *model.Order) error {
	return s.orderRepository.Update(ctx, order)
}

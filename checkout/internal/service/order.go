package service

import (
	"context"
	"ticket-system/checkout/internal/message"
	"ticket-system/checkout/internal/model"
	"ticket-system/checkout/internal/repository"
)

type orderService struct {
	orderRepository repository.OrderRepository
	orderMessenger  message.OrderMessenger
}

func NewOrderService(eventRepository repository.OrderRepository, orderMessenger message.OrderMessenger) OrderService {
	return &orderService{
		eventRepository,
		orderMessenger,
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

func (s orderService) CheckBookingsAreEnough(items []*model.Item, bookings []*model.Booking) error {
	stocksCounts := make(map[model.UUID]int)

	for _, i := range items {
		stocksCounts[i.StockId] += i.Count
	}

	for _, b := range bookings {
		if _, ok := stocksCounts[b.StockId]; !ok {
			return ErrNoStockBookingForItem
		}
		stocksCounts[b.StockId] -= b.Count
		if stocksCounts[b.StockId] == 0 {
			delete(stocksCounts, b.StockId)
		}
	}

	if len(stocksCounts) > 0 {
		return ErrBookedStocksAreNotEnoughForOrder
	}

	return nil
}

func (s orderService) ListOrders(ctx context.Context, userId model.UUID, limit int, offset int) ([]*model.Order, error) {
	filter := repository.FindOrdersByParams{
		UserId: model.NullUUID{Value: userId, Valid: true},
		Limit:  limit,
		Offset: offset,
	}
	return s.orderRepository.FindBy(ctx, filter)
}

func (s orderService) SendOrderMessage(ctx context.Context, o *model.Order) error {
	return s.orderMessenger.Send(ctx, o)
}

package service

import (
	"context"
	"ticket-system/booking/internal/model"
	"ticket-system/lib/query-engine/postgres"
)

type StockService interface {
	GetStock(ctx context.Context, uuid model.UUID) (*model.Stock, error)
	ListStocks(ctx context.Context, eventId model.UUID, ticketId *model.UUID) ([]*model.Stock, error)
	CreateStock(ctx context.Context, stock *model.Stock) error
	UpdateStock(ctx context.Context, stock *model.Stock) error
	AddBookedSeats(ctx context.Context, stock *model.Stock, count int) error
}

type BookingService interface {
	CreateBooking(ctx context.Context, booking *model.Booking) error
	ListBookings(ctx context.Context, stockId model.UUID, id model.UUID, userId model.UUID, expired bool) ([]*model.Booking, error)
}

type BusinessLogic struct {
	transactionManager *postgres.TransactionManager
	stockService       StockService
	bookingService     BookingService
}

func New(transactionManager *postgres.TransactionManager, stockService StockService, bookingService BookingService) BusinessLogic {
	return BusinessLogic{
		transactionManager,
		stockService,
		bookingService,
	}
}

func (s *BusinessLogic) GetStock(ctx context.Context, uuid model.UUID) (*model.Stock, error) {
	return s.stockService.GetStock(ctx, uuid)
}

func (s *BusinessLogic) ListStocks(ctx context.Context, eventId model.UUID, ticketId *model.UUID) ([]*model.Stock, error) {
	return s.stockService.ListStocks(ctx, eventId, ticketId)
}

func (s *BusinessLogic) CreateStock(ctx context.Context, stock *model.Stock) error {
	return s.stockService.CreateStock(ctx, stock)
}

func (s *BusinessLogic) GetStockForBooking(ctx context.Context, eventId model.UUID, ticketId model.UUID, seatsNeeded int) (*model.Stock, error) {
	stocks, err := s.ListStocks(ctx, eventId, &ticketId)

	if err != nil {
		return nil, err
	}

	if len(stocks) == 0 {
		return nil, ErrStockNotFound
	}

	stock := stocks[0]

	if stock.IsFullyBooked() {
		return nil, ErrStockIsFullyBooked
	}

	if stock.SeatsBooked+seatsNeeded > stock.SeatsTotal {
		return nil, ErrStockIsNotEnough
	}

	return stock, nil
}

func (s *BusinessLogic) CreateBooking(ctx context.Context, stock *model.Stock, userId model.UUID, orderId model.UUID, count int) (*model.Booking, error) {
	var booking *model.Booking

	err := s.transactionManager.RunRepeatableRead(ctx, func(ctxTX context.Context) error {
		err := s.stockService.AddBookedSeats(ctxTX, stock, count)
		if err != nil {
			return err
		}

		if stock.SeatsBooked > stock.SeatsTotal {
			return ErrStockIsNotEnough
		}

		booking = &model.Booking{
			StockId: stock.Id,
			UserId:  userId,
			OrderId: orderId,
			Count:   count,
		}

		err = s.bookingService.CreateBooking(ctxTX, booking)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		stock.SeatsBooked -= count
		return nil, err
	}

	return booking, nil
}

func (s *BusinessLogic) ListBookings(ctx context.Context, stockId model.UUID, orderId model.UUID, userId model.UUID, withExpired bool) ([]*model.Booking, error) {
	return s.bookingService.ListBookings(ctx, stockId, orderId, userId, withExpired)
}

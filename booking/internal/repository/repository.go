package repository

import (
	"context"
	"ticket-system/booking/internal/model"
)

type BookingRepository interface {
	Find(ctx context.Context, id model.UUID) (*model.Booking, error)
	FindBy(ctx context.Context, params FindBookingsByParams) ([]*model.Booking, error)
	Create(ctx context.Context, booking *model.Booking) error
}

type StockRepository interface {
	Find(ctx context.Context, id model.UUID) (*model.Stock, error)
	FindBy(ctx context.Context, filter FindStocksByParams) ([]*model.Stock, error)
	Create(ctx context.Context, stock *model.Stock) error
	Update(ctx context.Context, stock *model.Stock) error
	AddBookedSeats(ctx context.Context, s *model.Stock, quantity int) error
}

type NullUUID struct {
	String model.UUID
	Valid  bool
}

type FindBookingsByParams struct {
	StockId     model.UUID
	UserId      model.UUID
	OrderId     model.UUID
	WithExpired bool
}

type FindStocksByParams struct {
	EventId  model.UUID
	TicketId struct {
		Use bool
		Val model.UUID
	}
}

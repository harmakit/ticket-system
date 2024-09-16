package repository

import (
	"context"
	"ticket-system/booking/internal/model"
)

type BookingRepository interface {
	Find(ctx context.Context, id model.UUID) (*model.Booking, error)
	FindBy(ctx context.Context, params FindBookingsByParams) ([]*model.Booking, error)
	Create(ctx context.Context, booking *model.Booking) error
	BatchDelete(ctx context.Context, bookings []*model.Booking) error
}

type StockRepository interface {
	Find(ctx context.Context, id model.UUID) (*model.Stock, error)
	FindBy(ctx context.Context, filter FindStocksByParams) ([]*model.Stock, error)
	Create(ctx context.Context, stock *model.Stock) error
	Update(ctx context.Context, stock *model.Stock) error
	ModifyBookedSeats(ctx context.Context, s *model.Stock, quantity int) error
}

type NullUUID struct {
	Value model.UUID
	Valid bool
}

type FindBookingsByParams struct {
	StockId     NullUUID
	UserId      NullUUID
	OrderId     NullUUID
	OnlyExpired bool
	WithExpired bool
	Limit       uint64
	Offset      uint64
}

type FindStocksByParams struct {
	EventId  model.UUID
	TicketId NullUUID
}

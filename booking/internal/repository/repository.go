package repository

import (
	"context"
	"database/sql"
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
	Delete(ctx context.Context, id model.UUID) error
}

const NewUUID = "gen_random_uuid()"

type FindBookingsByParams struct {
	Ids         []model.UUID
	StockId     model.NullUUID
	UserId      model.NullUUID
	OrderId     model.NullUUID
	OnlyExpired bool
	WithExpired bool
	Limit       uint64
	Offset      uint64
}

type FindStocksByParams struct {
	EventId  model.UUID
	TicketId model.NullUUID
}

func NullString(s model.NullString) sql.NullString {
	return sql.NullString{
		String: s.Value,
		Valid:  s.Valid,
	}
}

func NullUUID(s model.NullUUID) sql.NullString {
	return sql.NullString{
		String: string(s.Value),
		Valid:  s.Valid,
	}
}

package service

import (
	"context"
	"ticket-system/booking/internal/model"
	"ticket-system/booking/internal/repository"
	"time"
)

type bookingService struct {
	bookingRepository repository.BookingRepository
}

const BookingTimeout = 15 * 60

func NewBookingService(bookingRepository repository.BookingRepository) BookingService {
	return &bookingService{
		bookingRepository,
	}
}

func (s bookingService) CreateBooking(ctx context.Context, booking *model.Booking) error {
	booking.CreatedAt = time.Now()
	booking.ExpiredAt = booking.CreatedAt.Add(BookingTimeout * time.Second)
	return s.bookingRepository.Create(ctx, booking)
}

func (s bookingService) ListBookings(ctx context.Context, stockId *model.UUID, orderId *model.UUID, userId *model.UUID, expired bool) ([]*model.Booking, error) {
	filter := repository.FindBookingsByParams{
		WithExpired: expired,
	}
	if stockId != nil {
		filter.StockId = model.NullUUID{Valid: true, Value: *stockId}
	}
	if orderId != nil {
		filter.OrderId = model.NullUUID{Valid: true, Value: *orderId}
	}
	if userId != nil {
		filter.UserId = model.NullUUID{Valid: true, Value: *userId}
	}
	return s.bookingRepository.FindBy(ctx, filter)
}

func (s bookingService) ListExpiredBookings(ctx context.Context, limit uint64, offset uint64) ([]*model.Booking, error) {
	filter := repository.FindBookingsByParams{
		OnlyExpired: true,
		Limit:       limit,
		Offset:      offset,
	}
	return s.bookingRepository.FindBy(ctx, filter)
}

func (s bookingService) DeleteBookings(ctx context.Context, bookings []*model.Booking) error {
	return s.bookingRepository.BatchDelete(ctx, bookings)
}

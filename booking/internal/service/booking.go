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

func (s bookingService) ListBookings(ctx context.Context, stockId model.UUID, orderId model.UUID, userId model.UUID, expired bool) ([]*model.Booking, error) {
	filter := repository.FindBookingsByParams{
		StockId:     stockId,
		OrderId:     orderId,
		UserId:      userId,
		WithExpired: expired,
	}
	return s.bookingRepository.FindBy(ctx, filter)
}

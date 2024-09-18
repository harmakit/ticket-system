package service

import (
	"errors"
)

var ErrStockNotFound = errors.New("stock not found")
var ErrStockIsFullyBooked = errors.New("stock is fully booked")
var ErrStockIsNotEnough = errors.New("stock is not enough")
var ErrStockIsNegative = errors.New("stock is negative")
var ErrBookingsMissing = errors.New("bookings missing")

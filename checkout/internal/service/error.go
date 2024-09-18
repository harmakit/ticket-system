package service

import "errors"

var ErrMultipleCartsReceived = errors.New("multiple carts received")
var ErrEmptyCart = errors.New("empty cart")
var ErrNoBookings = errors.New("no bookings")
var ErrNoStockBookingForItem = errors.New("no stock booking for item")
var ErrBookedStocksAreNotEnoughForOrder = errors.New("booked stocks are not enough for order")
var ErrOrderWrongStatus = errors.New("order wrong status")

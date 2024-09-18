package booking

import "errors"

var ErrGetBookings = errors.New("get bookings error")
var ErrDeleteStock = errors.New("delete stock error")
var ErrGetStocks = errors.New("get stocks error")
var ErrExpireBookings = errors.New("expire bookings error")
var ErrCreateBooking = errors.New("create booking error")
var ErrDeleteOrderBookings = errors.New("delete order bookings error")

var ErrNoStock = errors.New("no stock")
var ErrMultipleStocksReceived = errors.New("multiple stocks received")

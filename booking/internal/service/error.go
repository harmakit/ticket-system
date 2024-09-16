package service

import "fmt"

var ErrStockNotFound = fmt.Errorf("stock not found")
var ErrStockIsFullyBooked = fmt.Errorf("stock is fully booked")
var ErrStockIsNotEnough = fmt.Errorf("stock is not enough")
var ErrStockIsNegative = fmt.Errorf("stock is negative")

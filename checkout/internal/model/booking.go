package model

import "time"

type Booking struct {
	Id        UUID
	StockId   UUID
	UserId    UUID
	OrderId   UUID
	Count     int
	CreatedAt time.Time
	ExpiredAt time.Time
}

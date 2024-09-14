package schema

import (
	"time"
)

const BookingTable = "booking"

var BookingColumns = []string{"id", "stock_id", "user_id", "order_id", "count", "created_at", "expired_at"}

type Booking struct {
	Id        string
	StockId   string
	UserId    string
	OrderId   string
	Count     int
	CreatedAt time.Time
	ExpiredAt time.Time
}

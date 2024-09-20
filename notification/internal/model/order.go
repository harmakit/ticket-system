package model

type Order struct {
	Id     UUID
	UserId UUID
	Status OrderStatus
}

type OrderStatus string

const (
	StatusCreated   OrderStatus = "created"
	StatusPaid      OrderStatus = "paid"
	StatusFailed    OrderStatus = "failed"
	StatusCancelled OrderStatus = "cancelled"
)

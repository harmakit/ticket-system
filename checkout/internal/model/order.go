package model

type Order struct {
	Id     UUID
	UserId UUID
	Status Status
}

type Status string

const (
	StatusCreated   Status = "created"
	StatusPaid      Status = "paid"
	StatusFailed    Status = "failed"
	StatusCancelled Status = "cancelled"
)

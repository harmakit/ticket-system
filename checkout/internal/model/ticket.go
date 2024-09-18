package model

type Ticket struct {
	Id      UUID
	EventId UUID
	Name    string
	Price   float32
}

package model

type Item struct {
	Id       UUID
	OrderId  UUID
	StockId  UUID
	TicketId UUID
	Count    int
}

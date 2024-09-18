package model

type Cart struct {
	Id       UUID
	UserId   UUID
	TicketId UUID
	Count    int
}

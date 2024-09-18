package model

type Stock struct {
	Id          UUID
	EventId     UUID
	TicketId    UUID
	SeatsTotal  int
	SeatsBooked int
}

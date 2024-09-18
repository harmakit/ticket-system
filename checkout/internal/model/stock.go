package model

type Stock struct {
	Id          UUID
	EventId     UUID
	TicketId    UUID
	SeatsTotal  int
	SeatsBooked int
}

func (s *Stock) IsFullyBooked() bool {
	return s.SeatsTotal == s.SeatsBooked
}

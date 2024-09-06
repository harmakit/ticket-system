package schema

const TicketTable = "ticket"

var TicketColumns = []string{"id", "event_id", "name", "price"}

type Ticket struct {
	Id      string
	EventId string
	Name    string
	Price   float32
}

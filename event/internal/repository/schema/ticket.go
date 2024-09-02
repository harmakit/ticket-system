package schema

const TicketTable = "ticket"

var TicketColumns = []string{"event_id", "name", "price"}

type Ticket struct {
	EventId string
	Name    string
	Price   float32
}

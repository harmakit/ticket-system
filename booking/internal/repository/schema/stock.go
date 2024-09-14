package schema

const StockTable = "stock"

var StockColumns = []string{"id", "event_id", "ticket_id", "seats_total", "seats_booked"}

type Stock struct {
	Id          string
	EventId     string
	TicketId    string
	SeatsTotal  int
	SeatsBooked int
}

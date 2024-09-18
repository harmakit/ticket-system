package schema

const CartTable = "Cart"

var CartColumns = []string{"id", "user_id", "ticket_id", "count"}

type Cart struct {
	Id       string
	UserId   string
	TicketId string
	Count    int
}

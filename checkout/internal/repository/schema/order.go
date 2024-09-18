package schema

const OrderTable = "Order"

var OrderColumns = []string{"id", "user_id", "status"}

type Order struct {
	Id     string
	UserId string
	Status string
}

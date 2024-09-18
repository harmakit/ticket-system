package schema

const OrderTable = "\"order\""

var OrderColumns = []string{"id", "user_id", "status"}

type Order struct {
	Id     string
	UserId string
	Status string
}

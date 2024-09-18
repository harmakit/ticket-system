package schema

const ItemTable = "Item"

var ItemColumns = []string{"id", "order_id", "stock_id", "count"}

type Item struct {
	Id      string
	OrderId string
	StockId string
	Count   int
}

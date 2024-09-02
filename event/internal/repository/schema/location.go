package schema

const LocationTable = "location"

var LocationColumns = []string{"id", "name", "address", "lat", "lng"}

type Location struct {
	Id      string
	Name    string
	Address string
	Lat     float32
	Lng     float32
}

package schema

import (
	"database/sql"
	"time"
)

const EventTable = "event"

var EventColumns = []string{"id", "date", "duration", "name", "description", "location_id"}

type Event struct {
	Id          string
	Date        time.Time
	Duration    int
	Name        string
	Description sql.NullString
	LocationId  sql.NullString
}

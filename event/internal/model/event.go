package model

import (
	"database/sql"
	"time"
)

type Event struct {
	Id          UUID
	Date        time.Time
	Duration    int
	Name        string
	Description sql.NullString
	LocationId  sql.NullString
}

func (e *Event) IsOnline() bool {
	return !e.LocationId.Valid
}

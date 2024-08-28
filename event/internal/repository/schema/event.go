package model

import (
	"database/sql"
	"time"
)

type Event struct {
	id          string
	date        time.Time
	duration    int
	name        string
	description string
	location    sql.NullString
}

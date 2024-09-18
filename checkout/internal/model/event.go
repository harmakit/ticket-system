package model

import (
	"time"
)

type Event struct {
	Id          UUID
	Date        time.Time
	Duration    int
	Name        string
	Description NullString
}

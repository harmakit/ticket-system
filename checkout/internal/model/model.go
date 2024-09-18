package model

type UUID string

type NullUUID struct {
	Value UUID
	Valid bool
}

type NullString struct {
	Value string
	Valid bool
}

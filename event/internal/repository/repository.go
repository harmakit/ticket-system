package repository

import (
	"context"
	"ticket-system/event/internal/model"
)

type EventRepository interface {
	Find(ctx context.Context, id model.UUID) (*model.Event, error)
	FindBy(ctx context.Context, params FindEventsByParams) ([]*model.Event, error)
}

type TicketRepository interface {
	FindByEventId(ctx context.Context, eventId model.UUID) ([]*model.Ticket, error)
}

type LocationRepository interface {
	Find(ctx context.Context, id model.UUID) (*model.Location, error)
}

type NullUUID struct {
	String model.UUID
	Valid  bool
}

type FindEventsByParams struct {
	Offset     int
	Limit      int
	LocationId NullUUID
}

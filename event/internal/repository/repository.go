package repository

import (
	"context"
	"database/sql"
	"ticket-system/event/internal/model"
)

type EventRepository interface {
	Find(ctx context.Context, id string) (*model.Event, error)
	FindBy(ctx context.Context, params FindEventByParams) ([]*model.Event, error)
}

type TicketRepository interface {
	FindByEventId(ctx context.Context, eventId string) ([]*model.Ticket, error)
}

type LocationRepository interface {
	Find(ctx context.Context, id string) (*model.Location, error)
}

type FindEventByParams struct {
	Offset     int
	Limit      int
	LocationId sql.NullString
}

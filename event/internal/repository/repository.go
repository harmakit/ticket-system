package repository

import (
	"context"
	"database/sql"
	"ticket-system/event/internal/model"
)

type EventRepository interface {
	Find(ctx context.Context, id model.UUID) (*model.Event, error)
	FindBy(ctx context.Context, params FindEventsByParams) ([]*model.Event, error)
	Create(ctx context.Context, event *model.Event) error
}

type TicketRepository interface {
	Find(ctx context.Context, id model.UUID) (*model.Ticket, error)
	FindByEventId(ctx context.Context, eventId model.UUID) ([]*model.Ticket, error)
	Create(ctx context.Context, ticket *model.Ticket) error
}

type LocationRepository interface {
	Find(ctx context.Context, id model.UUID) (*model.Location, error)
}

type FindEventsByParams struct {
	Offset     int
	Limit      int
	LocationId model.NullUUID
}

func NullString(s model.NullString) sql.NullString {
	return sql.NullString{
		String: s.Value,
		Valid:  s.Valid,
	}
}

func NullUUID(s model.NullUUID) sql.NullString {
	return sql.NullString{
		String: string(s.Value),
		Valid:  s.Valid,
	}
}

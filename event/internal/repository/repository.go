package repository

import (
	"context"
	"ticket-system/event/internal/model"
)

type EventRepository interface {
	Find(ctx context.Context, id string) (*model.Event, error)
}

type TicketRepository interface {
	Find(ctx context.Context, id string) (*model.Ticket, error)
}

type LocationRepository interface {
	Find(ctx context.Context, id string) (*model.Location, error)
}

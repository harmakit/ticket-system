package service

import (
	"context"
	"database/sql"
	"ticket-system/event/internal/model"
	"ticket-system/event/internal/repository"
)

type eventService struct {
	eventRepository repository.EventRepository
}

type EventService interface {
	GetEvent(ctx context.Context, uuid model.UUID) (*model.Event, error)
	GetEvents(ctx context.Context, filter ListEventsFilter) ([]*model.Event, error)
}

func NewEventService(eventRepository repository.EventRepository) EventService {
	return &eventService{
		eventRepository,
	}
}

type ListEventsFilter struct {
	Offset     int
	Limit      int
	LocationId struct {
		Use bool
		Id  model.UUID
	}
}

func (s *eventService) GetEvent(ctx context.Context, uuid model.UUID) (*model.Event, error) {
	return s.eventRepository.Find(ctx, string(uuid))
}

func (s *eventService) GetEvents(ctx context.Context, filter ListEventsFilter) ([]*model.Event, error) {
	params := repository.FindEventByParams{
		Offset: filter.Offset,
		Limit:  filter.Limit,
		LocationId: sql.NullString{
			String: string(filter.LocationId.Id),
			Valid:  filter.LocationId.Use,
		},
	}
	return s.eventRepository.FindBy(ctx, params)
}

package service

import (
	"context"
	"github.com/pkg/errors"
	"ticket-system/event/internal/model"
	"ticket-system/event/internal/repository"
)

type eventService struct {
	eventRepository repository.EventRepository
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
		Val model.UUID
	}
}

func (s *eventService) GetEvent(ctx context.Context, uuid model.UUID) (*model.Event, error) {
	return s.eventRepository.Find(ctx, uuid)
}

func (s *eventService) GetEvents(ctx context.Context, filter ListEventsFilter) ([]*model.Event, error) {
	params := repository.FindEventsByParams{
		Offset: filter.Offset,
		Limit:  filter.Limit,
		LocationId: repository.NullUUID{
			String: filter.LocationId.Val,
			Valid:  filter.LocationId.Use,
		},
	}

	events, err := s.eventRepository.FindBy(ctx, params)
	if err != nil {
		if errors.Is(err, repository.ErrNoRows) {
			return nil, nil
		}
	}

	return events, err
}

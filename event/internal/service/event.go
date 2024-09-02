package service

import (
	"context"
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

func (s *eventService) GetEvent(ctx context.Context, id string) (*model.Event, error) {
	return s.eventRepository.Find(ctx, id)
}

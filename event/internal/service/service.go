package service

import (
	"context"
	"ticket-system/event/internal/model"
)

type BusinessLogic struct {
	eventService EventService
}

type EventService interface {
	GetEvent(ctx context.Context, id string) (*model.Event, error)
}

func New(eventService EventService) BusinessLogic {
	return BusinessLogic{
		eventService,
	}
}

func (s *BusinessLogic) GetEvent(ctx context.Context, id string) (*model.Event, error) {
	return s.eventService.GetEvent(ctx, id)
}

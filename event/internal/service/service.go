package service

import (
	"context"
	"ticket-system/event/internal/model"
)

type BusinessLogic struct {
	eventService    EventService
	locationService LocationService
	ticketService   TicketService
}

func New(eventService EventService, locationService LocationService, ticketService TicketService) BusinessLogic {
	return BusinessLogic{
		eventService,
		locationService,
		ticketService,
	}
}

func (s *BusinessLogic) GetEvent(ctx context.Context, uuid model.UUID) (*model.Event, error) {
	return s.eventService.GetEvent(ctx, uuid)
}

func (s *BusinessLogic) GetLocation(ctx context.Context, uuid model.UUID) (*model.Location, error) {
	return s.locationService.GetLocation(ctx, uuid)
}

func (s *BusinessLogic) GetTicketsForEvent(ctx context.Context, event *model.Event) ([]*model.Ticket, error) {
	return s.ticketService.GetTicketsForEvent(ctx, event)
}

func (s *BusinessLogic) ListEvents(ctx context.Context, limit, offset int, locationId *string) ([]*model.Event, []*model.Location, [][]*model.Ticket, error) {
	filter := ListEventsFilter{}
	filter.Limit = limit
	filter.Offset = offset
	if locationId != nil {
		filter.LocationId.Use = true
		filter.LocationId.Id = model.UUID(*locationId)
	}

	events, err := s.eventService.GetEvents(ctx, filter)
	if err != nil {
		return nil, nil, nil, err
	}

	locations := make([]*model.Location, len(events))
	tickets := make([][]*model.Ticket, len(events))

	for i, event := range events { // todo: worker pool
		if event.LocationId.Valid {
			locations[i], err = s.locationService.GetLocation(ctx, model.UUID(event.LocationId.String))
			if err != nil {
				return nil, nil, nil, err
			}
		}
		tickets[i], err = s.ticketService.GetTicketsForEvent(ctx, event)
		if err != nil {
			return nil, nil, nil, err
		}
	}

	return events, locations, tickets, err
}

package service

import (
	"context"
	"ticket-system/event/internal/model"
)

type EventService interface {
	GetEvent(ctx context.Context, uuid model.UUID) (*model.Event, error)
	GetEvents(ctx context.Context, filter ListEventsFilter) ([]*model.Event, error)
}

type LocationService interface {
	GetLocation(ctx context.Context, uuid model.UUID) (*model.Location, error)
	GetLocationsForEvents(ctx context.Context, events []*model.Event) ([]*model.Location, error)
}

type TicketService interface {
	GetTicketsForEvent(ctx context.Context, event *model.Event) ([]*model.Ticket, error)
	GetTicketsForEvents(ctx context.Context, events []*model.Event) ([][]*model.Ticket, error)
}

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
	filter := ListEventsFilter{
		Limit:  limit,
		Offset: offset,
	}
	if locationId != nil {
		filter.LocationId.Use = true
		filter.LocationId.Val = model.UUID(*locationId)
	}

	events, err := s.eventService.GetEvents(ctx, filter)
	if err != nil {
		return nil, nil, nil, err
	}

	locChan := make(chan []*model.Location, 1)
	ticketChan := make(chan [][]*model.Ticket, 1)
	errChan := make(chan error, 2)

	go func() {
		locations, err := s.locationService.GetLocationsForEvents(ctx, events)
		if err != nil {
			errChan <- err
			return
		}
		locChan <- locations
	}()

	go func() {
		tickets, err := s.ticketService.GetTicketsForEvents(ctx, events)
		if err != nil {
			errChan <- err
			return
		}
		ticketChan <- tickets
	}()

	var locations []*model.Location
	var tickets [][]*model.Ticket

	for i := 0; i < 2; i++ {
		select {
		case err := <-errChan:
			return nil, nil, nil, err
		case locs := <-locChan:
			locations = locs
		case tics := <-ticketChan:
			tickets = tics
		}
	}

	return events, locations, tickets, nil
}

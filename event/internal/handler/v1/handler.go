package v1

import (
	"context"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"ticket-system/event/internal/model"
	"ticket-system/event/internal/service"
	desc "ticket-system/event/pkg/v1/api"
)

type EventServiceImplementation struct {
	desc.UnsafeEventServiceServer
	bl service.BusinessLogic
}

func NewEventServiceImplementation(bl service.BusinessLogic) desc.EventServiceServer {
	return &EventServiceImplementation{bl: bl}
}

func (impl EventServiceImplementation) GetTicket(ctx context.Context, req *desc.GetTicketRequest) (*desc.GetTicketResponse, error) {
	var res desc.GetTicketResponse

	ticket, err := impl.bl.GetTicket(ctx, model.UUID(req.Id))
	if err != nil {
		return nil, err
	}

	res.Ticket = impl.bindModelToDescTicket(ticket)

	return &res, nil
}

func (impl EventServiceImplementation) GetEvent(ctx context.Context, req *desc.GetEventRequest) (*desc.GetEventResponse, error) {
	var res desc.GetEventResponse

	event, err := impl.bl.GetEvent(ctx, model.UUID(req.GetId()))
	if err != nil {
		return nil, err
	}

	var location *model.Location
	if event.LocationId.Valid {
		location, err = impl.bl.GetLocation(ctx, event.LocationId.Value)
		if err != nil {
			return nil, err
		}
	}

	var tickets []*model.Ticket
	tickets, err = impl.bl.GetTicketsForEvent(ctx, event)
	if err != nil {
		return nil, err
	}

	res.Event = impl.bindModelToDescEvent(event, location, tickets)

	return &res, nil
}

func (impl EventServiceImplementation) ListEvents(ctx context.Context, req *desc.ListEventsRequest) (*desc.ListEventsResponse, error) {
	var res desc.ListEventsResponse

	events, locations, tickets, err := impl.bl.ListEvents(ctx, int(req.Limit), int(req.Offset), req.LocationId)
	if err != nil {
		return nil, err
	}

	res.Events = make([]*desc.Event, len(events))

	for i := range len(res.Events) {
		res.Events[i] = impl.bindModelToDescEvent(events[i], locations[i], tickets[i])
	}

	return &res, nil
}

func (impl EventServiceImplementation) CreateEvent(ctx context.Context, req *desc.CreateEventRequest) (*desc.CreateEventResponse, error) {
	var res desc.CreateEventResponse

	event := &model.Event{
		Name:     req.Event.Name,
		Date:     req.Event.Date.AsTime(),
		Duration: int(req.Event.Duration),
		Description: model.NullString{
			Value: req.Event.Description,
			Valid: true,
		},
	}

	if req.Event.LocationId != nil {
		event.LocationId = model.NullUUID{
			Value: model.UUID(*req.Event.LocationId),
			Valid: true,
		}
	}

	tickets := make([]*model.Ticket, len(req.Event.Tickets))
	for i, ticket := range req.Event.Tickets {
		tickets[i] = &model.Ticket{
			Name:  ticket.Name,
			Price: ticket.Price,
		}
	}

	seats := make([]int, len(req.Event.Tickets))
	for i, ticket := range req.Event.Tickets {
		seats[i] = int(ticket.Seats)
	}

	var location *model.Location
	var err error

	if event.LocationId.Valid {
		location, err = impl.bl.GetLocation(ctx, event.LocationId.Value)
		if err != nil {
			return nil, err
		}
	}

	err = impl.bl.CreateEvent(ctx, event, tickets, seats)
	if err != nil {
		return nil, err
	}

	res.Event = impl.bindModelToDescEvent(event, location, tickets)

	return &res, nil
}

func (impl EventServiceImplementation) bindModelToDescEvent(event *model.Event, location *model.Location, tickets []*model.Ticket) *desc.Event {
	if event == nil {
		return nil
	}

	var dl *desc.Location
	if location != nil {
		dl = &desc.Location{
			Id:      string(location.Id),
			Name:    location.Name,
			Address: location.Address,
			Lat:     location.Lat,
			Lng:     location.Lng,
		}
	}

	dts := make([]*desc.Ticket, len(tickets))
	for i, ticket := range tickets {
		dts[i] = impl.bindModelToDescTicket(ticket)
	}

	description := &wrapperspb.StringValue{}
	if event.Description.Valid {
		description.Value = event.Description.Value
	}

	de := &desc.Event{
		Id:          string(event.Id),
		Date:        timestamppb.New(event.Date),
		Duration:    int32(event.Duration),
		Name:        event.Name,
		Description: description,
		Location:    dl,
		Tickets:     dts,
	}

	return de
}

func (impl EventServiceImplementation) bindModelToDescTicket(ticket *model.Ticket) *desc.Ticket {
	return &desc.Ticket{
		Id:      string(ticket.Id),
		EventId: string(ticket.EventId),
		Price:   ticket.Price,
		Name:    ticket.Name,
	}
}

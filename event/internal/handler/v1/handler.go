package v1

import (
	"context"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"ticket-system/event/internal/model"
	"ticket-system/event/internal/service"
	desc "ticket-system/event/pkg/v1"
)

type EventServiceImplementation struct {
	desc.UnsafeEventServiceServer
	bl service.BusinessLogic
}

func NewEventServiceImplementation(bl service.BusinessLogic) desc.EventServiceServer {
	return &EventServiceImplementation{bl: bl}
}

func (impl EventServiceImplementation) GetEvent(ctx context.Context, req *desc.GetEventRequest) (*desc.GetEventResponse, error) {
	var res desc.GetEventResponse

	event, err := impl.bl.GetEvent(ctx, model.UUID(req.GetId()))
	if err != nil {
		return nil, err
	}

	var location *model.Location
	if event.LocationId.Valid {
		location, err = impl.bl.GetLocation(ctx, model.UUID(event.LocationId.String))
		if err != nil {
			return nil, err
		}
	}

	var tickets []*model.Ticket
	tickets, err = impl.bl.GetTicketsForEvent(ctx, event)
	if err != nil {
		return nil, err
	}

	res.Event = impl.bindModelsToDescEvent(event, location, tickets)

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
		res.Events[i] = impl.bindModelsToDescEvent(events[i], locations[i], tickets[i])
	}

	return &res, nil
}

func (impl EventServiceImplementation) bindModelsToDescEvent(event *model.Event, location *model.Location, tickets []*model.Ticket) *desc.Event {

	var description *wrapperspb.StringValue
	if event.Description.Valid {
		description.Value = event.Description.String
	}

	dl := &desc.Location{
		Id:      string(location.Id),
		Name:    location.Name,
		Address: location.Address,
		Lat:     location.Lat,
		Lng:     location.Lng,
	}

	dts := make([]*desc.Ticket, len(tickets))
	for i, ticket := range tickets {
		dts[i] = &desc.Ticket{
			Id:      string(ticket.Id),
			EventId: string(ticket.EventId),
			Price:   ticket.Price,
			Name:    ticket.Name,
		}
	}

	de := &desc.Event{
		Id:          string(event.Id),
		Date:        timestamppb.New(event.Date),
		Duration:    int32(event.Duration),
		Name:        event.Name,
		Description: description,
		Location:    dl,
		Tickets:     nil,
	}

	return de
}

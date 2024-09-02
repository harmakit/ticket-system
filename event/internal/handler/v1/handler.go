package v1

import (
	"context"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
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

func (i EventServiceImplementation) GetEvent(context context.Context, req *desc.GetEventRequest) (*desc.GetEventResponse, error) {
	var res desc.GetEventResponse

	event, err := i.bl.GetEvent(context, req.GetId())
	if err != nil {
		return nil, err
	}

	var description *wrapperspb.StringValue
	if event.Description.Valid {
		description.Value = event.Description.String
	}

	res.Event = &desc.Event{
		Id:          string(event.Id),
		Date:        timestamppb.New(event.Date),
		Duration:    int32(event.Duration),
		Name:        event.Name,
		Description: description,
		Location:    nil,
		Tickets:     nil,
	}

	return &res, nil
}

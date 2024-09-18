package event

import (
	"context"
	"go.uber.org/zap"
	"ticket-system/checkout/internal/config"
	"ticket-system/checkout/internal/model"
	eventAPI "ticket-system/event/pkg/v1"
	grpcClient "ticket-system/lib/client"
)

type Service interface {
	GetTicket(ctx context.Context, id model.UUID) (*model.Ticket, error)
	GetEvent(ctx context.Context, id model.UUID) (*model.Event, error)
}

type service struct {
	client Client
}

func NewService() Service {
	eventClient := newClient(grpcClient.GetGRPCConn(config.Data.Services.Event, zap.NewNop()))
	return &service{eventClient}
}

func (s *service) GetTicket(ctx context.Context, id model.UUID) (*model.Ticket, error) {
	req := &eventAPI.GetTicketRequest{Id: string(id)}
	res, err := s.client.GetTicket(ctx, req)
	if err != nil {
		return nil, err
	}

	ticket := s.bindAPITicketToModel(res.Ticket)
	return ticket, nil
}

func (s *service) GetEvent(ctx context.Context, id model.UUID) (*model.Event, error) {
	req := &eventAPI.GetEventRequest{Id: string(id)}
	res, err := s.client.GetEvent(ctx, req)
	if err != nil {
		return nil, err
	}

	event := s.bindAPIEventToModel(res.Event)
	return event, nil
}

func (s *service) bindAPITicketToModel(t *eventAPI.Ticket) *model.Ticket {
	return &model.Ticket{
		Id:      model.UUID(t.Id),
		EventId: model.UUID(t.EventId),
		Name:    t.Name,
		Price:   t.Price,
	}
}

func (s *service) bindAPIEventToModel(e *eventAPI.Event) *model.Event {
	event := &model.Event{
		Id:       model.UUID(e.Id),
		Date:     e.Date.AsTime(),
		Duration: int(e.Duration),
		Name:     e.Name,
	}

	var description model.NullString
	if e.Description != nil {
		description = model.NullString{Valid: true, Value: e.Description.Value}
	}
	event.Description = description

	return event
}

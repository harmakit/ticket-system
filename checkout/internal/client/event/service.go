package event

import (
	"context"
	"go.uber.org/zap"
	"ticket-system/checkout/internal/config"
	"ticket-system/checkout/internal/model"
	"ticket-system/event/pkg/v1/api"
	grpcClient "ticket-system/lib/client"
)

type Service interface {
	GetTicket(ctx context.Context, id model.UUID) (*model.Ticket, error)
}

type service struct {
	client Client
}

func NewService() Service {
	eventClient := newClient(grpcClient.GetGRPCConn(config.Data.Services.Event, zap.NewNop()))
	return &service{eventClient}
}

func (s *service) GetTicket(ctx context.Context, id model.UUID) (*model.Ticket, error) {
	req := &api.GetTicketRequest{Id: string(id)}
	res, err := s.client.GetTicket(ctx, req)
	if err != nil {
		return nil, err
	}

	ticket := s.bindAPITicketToModel(res.Ticket)
	return ticket, nil
}

func (s *service) bindAPITicketToModel(t *api.Ticket) *model.Ticket {
	return &model.Ticket{
		Id:      model.UUID(t.Id),
		EventId: model.UUID(t.EventId),
		Name:    t.Name,
		Price:   t.Price,
	}
}

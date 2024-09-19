package booking

import (
	"context"
	"go.uber.org/zap"
	"ticket-system/booking/pkg/v1/api"
	"ticket-system/event/internal/config"
	"ticket-system/event/internal/model"
	grpcClient "ticket-system/lib/client"
)

type Service interface {
	CreateStock(ctx context.Context, eventId model.UUID, ticketId model.UUID, seats int) (model.UUID, error)
	DeleteStock(ctx context.Context, id model.UUID) error
}

type service struct {
	client Client
}

func NewService() Service {
	bookingClient := newClient(grpcClient.GetGRPCConn(config.Data.Services.Booking, zap.NewNop()))
	return &service{bookingClient}
}

func (s *service) CreateStock(ctx context.Context, eventId model.UUID, ticketId model.UUID, seats int) (model.UUID, error) {
	req := &api.CreateStockRequest{
		EventId:  string(eventId),
		TicketId: string(ticketId),
		Seats:    int32(seats),
	}
	res, err := s.client.CreateStock(ctx, req)
	if err != nil {
		return "", err
	}
	return model.UUID(res.Id), nil
}

func (s *service) DeleteStock(ctx context.Context, id model.UUID) error {
	req := &api.DeleteStockRequest{Id: string(id)}
	_, err := s.client.DeleteStock(ctx, req)
	return err
}

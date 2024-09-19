package booking

import (
	"context"
	"go.uber.org/zap"
	"ticket-system/booking/pkg/v1/api"
	"ticket-system/checkout/internal/config"
	"ticket-system/checkout/internal/model"
	grpcClient "ticket-system/lib/client"
)

type Service interface {
	GetTicketStock(ctx context.Context, eventId model.UUID, ticketId model.UUID) (*model.Stock, error)
	ExpireBookings(ctx context.Context, bookingsIds []model.UUID) error
	CreateBooking(ctx context.Context, ticket *model.Ticket, order *model.Order, cart *model.Cart, userId model.UUID) (model.UUID, error)
	DeleteOrderBookings(ctx context.Context, orderId model.UUID, userId model.UUID, bookings []*model.Booking) error
	GetOrderBookings(ctx context.Context, order *model.Order, items []*model.Item, tickets []*model.Ticket) ([]*model.Booking, error)
}

type service struct {
	client Client
}

func NewService() Service {
	bookingClient := newClient(grpcClient.GetGRPCConn(config.Data.Services.Booking, zap.NewNop()))
	return &service{bookingClient}
}

func (s *service) GetTicketStock(ctx context.Context, eventId model.UUID, ticketId model.UUID) (*model.Stock, error) {
	ticketIdVal := string(ticketId)
	req := &api.GetStocksRequest{EventId: string(eventId), TicketId: &ticketIdVal}
	res, err := s.client.GetStocks(ctx, req)
	if err != nil {
		return nil, err
	}

	stocks := make([]*model.Stock, 0, len(res.Stocks))
	for _, apiStock := range res.Stocks {
		stocks = append(stocks, s.bindAPIStockToModel(apiStock))
	}

	if len(stocks) == 0 {
		return nil, ErrNoStock
	}
	if len(stocks) > 1 {
		return nil, ErrMultipleStocksReceived
	}

	return stocks[0], nil
}

func (s *service) ExpireBookings(ctx context.Context, bookingsIds []model.UUID) error {
	ids := make([]string, 0, len(bookingsIds))
	for _, id := range bookingsIds {
		ids = append(ids, string(id))
	}

	req := &api.ExpireBookingsRequest{Ids: ids}
	_, err := s.client.ExpireBookings(ctx, req)
	return err
}

func (s *service) CreateBooking(ctx context.Context, ticket *model.Ticket, order *model.Order, cart *model.Cart, userId model.UUID) (model.UUID, error) {
	ticketId := ticket.Id
	eventId := ticket.EventId
	orderId := order.Id
	count := cart.Count

	req := &api.CreateBookingRequest{
		EventId:  string(eventId),
		TicketId: string(ticketId),
		UserId:   string(userId),
		OrderId:  string(orderId),
		Count:    int32(count),
	}

	res, err := s.client.CreateBooking(ctx, req)
	if err != nil {
		return "", err
	}
	return model.UUID(res.Id), nil
}

func (s *service) DeleteOrderBookings(ctx context.Context, orderId model.UUID, userId model.UUID, bookings []*model.Booking) error {
	ids := make([]string, len(bookings))
	for i, b := range bookings {
		ids[i] = string(b.Id)
	}

	req := &api.DeleteOrderBookingsRequest{
		OrderId: string(orderId),
		UserId:  string(userId),
		Ids:     ids,
	}
	_, err := s.client.DeleteOrderBookings(ctx, req)
	return err
}

func (s *service) GetOrderBookings(ctx context.Context, order *model.Order, items []*model.Item, tickets []*model.Ticket) ([]*model.Booking, error) {
	if len(items) != len(tickets) {
		return nil, ErrMismatchedItemsAndTickets
	}

	var bookings []*model.Booking
	for i, item := range items {
		req := &api.GetOrderBookingsRequest{
			EventId:     string(tickets[i].EventId),
			TicketId:    string(item.TicketId),
			OrderId:     string(order.Id),
			UserId:      string(order.UserId),
			WithExpired: false,
		}
		res, err := s.client.GetBookings(ctx, req)
		if err != nil {
			return nil, err
		}

		for _, apiBooking := range res.Bookings {
			bookings = append(bookings, s.bindAPIBookingToModel(apiBooking))
		}
	}

	return bookings, nil
}

func (s *service) bindAPIStockToModel(stock *api.Stock) *model.Stock {
	return &model.Stock{
		Id:          model.UUID(stock.Id),
		EventId:     model.UUID(stock.EventId),
		TicketId:    model.UUID(stock.TicketId),
		SeatsTotal:  int(stock.SeatsTotal),
		SeatsBooked: int(stock.SeatsBooked),
	}
}

func (s *service) bindAPIBookingToModel(booking *api.Booking) *model.Booking {
	return &model.Booking{
		Id:        model.UUID(booking.Id),
		StockId:   model.UUID(booking.StockId),
		UserId:    model.UUID(booking.UserId),
		OrderId:   model.UUID(booking.OrderId),
		Count:     int(booking.Count),
		CreatedAt: booking.CreatedAt.AsTime(),
		ExpiredAt: booking.ExpiredAt.AsTime(),
	}
}

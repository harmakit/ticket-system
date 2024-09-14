package v1

import (
	"context"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/timestamppb"
	"ticket-system/booking/internal/model"
	"ticket-system/booking/internal/service"
	desc "ticket-system/booking/pkg/v1"
)

type BookingServiceImplementation struct {
	desc.UnsafeBookingServiceServer
	bl service.BusinessLogic
}

func NewBookingServiceImplementation(bl service.BusinessLogic) desc.BookingServiceServer {
	return &BookingServiceImplementation{bl: bl}
}

func (impl BookingServiceImplementation) GetStock(ctx context.Context, req *desc.GetStockRequest) (*desc.GetStockResponse, error) {
	var res desc.GetStockResponse

	stock, err := impl.bl.GetStock(ctx, model.UUID(req.Id))
	if err != nil {
		return nil, err
	}

	res.Stock = impl.bindModelToDescStock(stock)

	return &res, nil
}

func (impl BookingServiceImplementation) GetStocks(ctx context.Context, req *desc.GetStocksRequest) (*desc.GetStocksResponse, error) {
	var res desc.GetStocksResponse

	var ticketId *model.UUID
	if req.TicketId != nil {
		ticketIdVal := model.UUID(*req.TicketId)
		ticketId = &ticketIdVal
	}

	stocks, err := impl.bl.ListStocks(ctx, model.UUID(req.EventId), ticketId)
	if err != nil {
		return nil, err
	}

	for _, stock := range stocks {
		res.Stocks = append(res.Stocks, impl.bindModelToDescStock(stock))
	}

	return &res, nil
}

func (impl BookingServiceImplementation) CreateStock(ctx context.Context, req *desc.CreateStockRequest) (*desc.CreateStockResponse, error) {
	var res desc.CreateStockResponse

	stock := &model.Stock{
		EventId:     model.UUID(req.EventId),
		TicketId:    model.UUID(req.TicketId),
		SeatsTotal:  int(req.Seats),
		SeatsBooked: 0,
	}

	err := impl.bl.CreateStock(ctx, stock)
	if err != nil {
		return nil, err
	}

	res.Id = string(stock.Id)

	return &res, nil
}

func (impl BookingServiceImplementation) GetBookings(ctx context.Context, req *desc.GetOrderBookingsRequest) (*desc.GetOrderBookingsResponse, error) {
	var res desc.GetOrderBookingsResponse

	ticketId := model.UUID(req.TicketId)
	stocks, err := impl.bl.ListStocks(ctx, model.UUID(req.EventId), &ticketId)
	if err != nil {
		return nil, err
	}

	if len(stocks) == 0 {
		return nil, nil
	}

	stock := stocks[0]
	bookings, err := impl.bl.ListBookings(ctx, stock.Id, model.UUID(req.OrderId), model.UUID(req.UserId), req.WithExpired)
	if err != nil {
		return nil, err
	}

	for _, booking := range bookings {
		res.Bookings = append(res.Bookings, impl.bindModelToDescBooking(booking))
	}

	return &res, nil
}

func (impl BookingServiceImplementation) CreateBooking(ctx context.Context, req *desc.CreateBookingRequest) (*desc.CreateBookingResponse, error) {
	var res desc.CreateBookingResponse

	stock, err := impl.bl.GetStockForBooking(ctx, model.UUID(req.EventId), model.UUID(req.TicketId), int(req.Count))
	if err != nil {
		return nil, err
	}
	if stock == nil {
		return nil, errors.WithMessage(ErrRuntimeError, "bl.GetStockForBooking returned nil stock")
	}

	booking, err := impl.bl.CreateBooking(ctx, stock, model.UUID(req.UserId), model.UUID(req.OrderId), int(req.Count))
	if err != nil {
		return nil, err
	}
	if booking == nil {
		return nil, errors.WithMessage(ErrRuntimeError, "bl.CreateBooking returned nil booking")
	}

	res.Id = string(booking.Id)
	res.CreatedAt = timestamppb.New(booking.CreatedAt)
	res.ExpiredAt = timestamppb.New(booking.ExpiredAt)

	return &res, nil
}

func (impl BookingServiceImplementation) bindModelToDescStock(s *model.Stock) *desc.Stock {
	if s == nil {
		return nil
	}

	var ds *desc.Stock

	ds = &desc.Stock{
		Id:          string(s.Id),
		EventId:     string(s.EventId),
		TicketId:    string(s.TicketId),
		SeatsTotal:  int32(s.SeatsTotal),
		SeatsBooked: int32(s.SeatsBooked),
	}

	return ds
}

func (impl BookingServiceImplementation) bindModelToDescBooking(b *model.Booking) *desc.Booking {
	if b == nil {
		return nil
	}

	var db *desc.Booking

	db = &desc.Booking{
		Id:        string(b.Id),
		StockId:   string(b.StockId),
		UserId:    string(b.UserId),
		OrderId:   string(b.OrderId),
		Count:     int32(b.Count),
		CreatedAt: timestamppb.New(b.CreatedAt),
		ExpiredAt: timestamppb.New(b.ExpiredAt),
	}

	return db
}

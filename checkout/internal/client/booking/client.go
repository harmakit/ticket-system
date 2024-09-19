package booking

import (
	"context"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"ticket-system/booking/pkg/v1/api"
)

type Client interface {
	Close() error
	GetStocks(ctx context.Context, req *api.GetStocksRequest) (*api.GetStocksResponse, error)
	ExpireBookings(ctx context.Context, req *api.ExpireBookingsRequest) (*emptypb.Empty, error)
	CreateBooking(ctx context.Context, req *api.CreateBookingRequest) (*api.CreateBookingResponse, error)
	DeleteOrderBookings(ctx context.Context, req *api.DeleteOrderBookingsRequest) (*emptypb.Empty, error)
	GetBookings(ctx context.Context, req *api.GetOrderBookingsRequest) (*api.GetOrderBookingsResponse, error)
}

type client struct {
	api.BookingServiceClient
	conn *grpc.ClientConn
}

func newClient(conn *grpc.ClientConn) Client {
	return &client{
		BookingServiceClient: api.NewBookingServiceClient(conn),
		conn:                 conn,
	}
}

func (c *client) Close() error {
	return c.conn.Close()
}

func (c *client) GetStocks(ctx context.Context, req *api.GetStocksRequest) (*api.GetStocksResponse, error) {
	res, err := c.BookingServiceClient.GetStocks(ctx, req)
	err = errors.Wrap(err, ErrGetStocks.Error())
	return res, err
}

func (c *client) ExpireBookings(ctx context.Context, req *api.ExpireBookingsRequest) (*emptypb.Empty, error) {
	res, err := c.BookingServiceClient.ExpireBookings(ctx, req)
	err = errors.Wrap(err, ErrExpireBookings.Error())
	return res, err
}

func (c *client) CreateBooking(ctx context.Context, req *api.CreateBookingRequest) (*api.CreateBookingResponse, error) {
	res, err := c.BookingServiceClient.CreateBooking(ctx, req)
	err = errors.Wrap(err, ErrCreateBooking.Error())
	return res, err
}

func (c *client) DeleteOrderBookings(ctx context.Context, req *api.DeleteOrderBookingsRequest) (*emptypb.Empty, error) {
	res, err := c.BookingServiceClient.DeleteOrderBookings(ctx, req)
	err = errors.Wrap(err, ErrDeleteOrderBookings.Error())
	return res, err
}

func (c *client) GetBookings(ctx context.Context, req *api.GetOrderBookingsRequest) (*api.GetOrderBookingsResponse, error) {
	res, err := c.BookingServiceClient.GetBookings(ctx, req)
	err = errors.Wrap(err, ErrGetBookings.Error())
	return res, err
}

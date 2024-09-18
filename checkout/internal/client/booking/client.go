package booking

import (
	"context"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	bookingAPI "ticket-system/booking/pkg/v1"
)

type Client interface {
	Close() error
	GetBookings(ctx context.Context, req *bookingAPI.GetOrderBookingsRequest) (*bookingAPI.GetOrderBookingsResponse, error)
	DeleteStock(ctx context.Context, req *bookingAPI.DeleteStockRequest) (*emptypb.Empty, error)
	GetStocks(ctx context.Context, req *bookingAPI.GetStocksRequest) (*bookingAPI.GetStocksResponse, error)
	ExpireBookings(ctx context.Context, req *bookingAPI.ExpireBookingsRequest) (*emptypb.Empty, error)
	CreateBooking(ctx context.Context, req *bookingAPI.CreateBookingRequest) (*bookingAPI.CreateBookingResponse, error)
	DeleteOrderBookings(ctx context.Context, req *bookingAPI.DeleteOrderBookingsRequest) (*emptypb.Empty, error)
}

type client struct {
	bookingAPI.BookingServiceClient
	conn *grpc.ClientConn
}

func newClient(conn *grpc.ClientConn) Client {
	return &client{
		BookingServiceClient: bookingAPI.NewBookingServiceClient(conn),
		conn:                 conn,
	}
}

func (c *client) Close() error {
	return c.conn.Close()
}

func (c *client) GetBookings(ctx context.Context, req *bookingAPI.GetOrderBookingsRequest) (*bookingAPI.GetOrderBookingsResponse, error) {
	res, err := c.BookingServiceClient.GetBookings(ctx, req)
	err = errors.Wrap(err, ErrGetBookings.Error())
	return res, err
}

func (c *client) DeleteStock(ctx context.Context, req *bookingAPI.DeleteStockRequest) (*emptypb.Empty, error) {
	res, err := c.BookingServiceClient.DeleteStock(ctx, req)
	err = errors.Wrap(err, ErrDeleteStock.Error())
	return res, err
}

func (c *client) GetStocks(ctx context.Context, req *bookingAPI.GetStocksRequest) (*bookingAPI.GetStocksResponse, error) {
	res, err := c.BookingServiceClient.GetStocks(ctx, req)
	err = errors.Wrap(err, ErrGetStocks.Error())
	return res, err
}

func (c *client) ExpireBookings(ctx context.Context, req *bookingAPI.ExpireBookingsRequest) (*emptypb.Empty, error) {
	res, err := c.BookingServiceClient.ExpireBookings(ctx, req)
	err = errors.Wrap(err, ErrExpireBookings.Error())
	return res, err
}

func (c *client) CreateBooking(ctx context.Context, req *bookingAPI.CreateBookingRequest) (*bookingAPI.CreateBookingResponse, error) {
	res, err := c.BookingServiceClient.CreateBooking(ctx, req)
	err = errors.Wrap(err, ErrCreateBooking.Error())
	return res, err
}

func (c *client) DeleteOrderBookings(ctx context.Context, req *bookingAPI.DeleteOrderBookingsRequest) (*emptypb.Empty, error) {
	res, err := c.BookingServiceClient.DeleteOrderBookings(ctx, req)
	err = errors.Wrap(err, ErrDeleteOrderBookings.Error())
	return res, err
}

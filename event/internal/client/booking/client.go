package booking

import (
	"context"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	bookingAPI "ticket-system/booking/pkg/v1"

	"google.golang.org/protobuf/types/known/emptypb"
)

type Client interface {
	Close() error
	CreateStock(ctx context.Context, req *bookingAPI.CreateStockRequest) (*bookingAPI.CreateStockResponse, error)
	DeleteStock(ctx context.Context, req *bookingAPI.DeleteStockRequest) (*emptypb.Empty, error)
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

func (c *client) CreateStock(ctx context.Context, req *bookingAPI.CreateStockRequest) (*bookingAPI.CreateStockResponse, error) {
	res, err := c.BookingServiceClient.CreateStock(ctx, req)
	err = errors.Wrap(err, ErrCreateStock.Error())
	return res, err
}

func (c *client) DeleteStock(ctx context.Context, req *bookingAPI.DeleteStockRequest) (*emptypb.Empty, error) {
	res, err := c.BookingServiceClient.DeleteStock(ctx, req)
	err = errors.Wrap(err, ErrDeleteStock.Error())
	return res, err
}

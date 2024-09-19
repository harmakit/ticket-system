package booking

import (
	"context"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"ticket-system/booking/pkg/v1/api"

	"google.golang.org/protobuf/types/known/emptypb"
)

type Client interface {
	Close() error
	CreateStock(ctx context.Context, req *api.CreateStockRequest) (*api.CreateStockResponse, error)
	DeleteStock(ctx context.Context, req *api.DeleteStockRequest) (*emptypb.Empty, error)
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

func (c *client) CreateStock(ctx context.Context, req *api.CreateStockRequest) (*api.CreateStockResponse, error) {
	res, err := c.BookingServiceClient.CreateStock(ctx, req)
	err = errors.Wrap(err, ErrCreateStock.Error())
	return res, err
}

func (c *client) DeleteStock(ctx context.Context, req *api.DeleteStockRequest) (*emptypb.Empty, error) {
	res, err := c.BookingServiceClient.DeleteStock(ctx, req)
	err = errors.Wrap(err, ErrDeleteStock.Error())
	return res, err
}

package event

import (
	"context"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"ticket-system/event/pkg/v1/api"
)

type Client interface {
	Close() error
	GetTicket(ctx context.Context, req *api.GetTicketRequest) (*api.GetTicketResponse, error)
}

type client struct {
	api.EventServiceClient
	conn *grpc.ClientConn
}

func newClient(conn *grpc.ClientConn) Client {
	return &client{
		EventServiceClient: api.NewEventServiceClient(conn),
		conn:               conn,
	}
}

func (c *client) Close() error {
	return c.conn.Close()
}

func (c *client) GetTicket(ctx context.Context, req *api.GetTicketRequest) (*api.GetTicketResponse, error) {
	res, err := c.EventServiceClient.GetTicket(ctx, req)
	err = errors.Wrap(err, ErrGetTicket.Error())
	return res, err
}

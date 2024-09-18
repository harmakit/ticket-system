package event

import (
	"context"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	eventAPI "ticket-system/event/pkg/v1"
)

type Client interface {
	Close() error
	GetTicket(ctx context.Context, req *eventAPI.GetTicketRequest) (*eventAPI.GetTicketResponse, error)
}

type client struct {
	eventAPI.EventServiceClient
	conn *grpc.ClientConn
}

func newClient(conn *grpc.ClientConn) Client {
	return &client{
		EventServiceClient: eventAPI.NewEventServiceClient(conn),
		conn:               conn,
	}
}

func (c *client) Close() error {
	return c.conn.Close()
}

func (c *client) GetTicket(ctx context.Context, req *eventAPI.GetTicketRequest) (*eventAPI.GetTicketResponse, error) {
	res, err := c.EventServiceClient.GetTicket(ctx, req)
	err = errors.Wrap(err, ErrGetTicket.Error())
	return res, err
}

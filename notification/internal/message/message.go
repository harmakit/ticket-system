package message

import (
	"context"
	desc "ticket-system/checkout/pkg/v1/message"
)

type Handler interface {
	HandleOrderMessage(ctx context.Context, m *desc.OrderMessage) error
}

type Messenger interface {
	Close() error
	Consume(ctx context.Context, handleErr func(err error))
}

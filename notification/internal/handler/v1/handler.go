package v1

import (
	"context"
	desc "ticket-system/checkout/pkg/v1/message"
	"ticket-system/notification/internal/message"
	"ticket-system/notification/internal/model"
	"ticket-system/notification/internal/service"
)

type MessageHandlerImplementation struct {
	bl service.BusinessLogic
}

func NewMessageHandlerImplementation(bl service.BusinessLogic) message.Handler {
	return &MessageHandlerImplementation{bl: bl}
}

func (impl *MessageHandlerImplementation) HandleOrderMessage(ctx context.Context, m *desc.OrderMessage) error {
	var status model.OrderStatus
	switch m.Status {
	case desc.OrderStatus_CREATED:
		status = model.StatusCreated
	case desc.OrderStatus_PAID:
		status = model.StatusPaid
	case desc.OrderStatus_FAILED:
		status = model.StatusFailed
	case desc.OrderStatus_CANCELLED:
		status = model.StatusCancelled
	}
	if status == "" {
		return ErrUnknownOrderStatus
	}

	order := &model.Order{
		Id:     model.UUID(m.Id),
		UserId: model.UUID(m.UserId),
		Status: status,
	}
	return impl.bl.NotifyOrderChanges(order)
}

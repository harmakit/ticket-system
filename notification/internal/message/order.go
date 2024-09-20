package message

import (
	"context"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	desc "ticket-system/checkout/pkg/v1/message"
	"ticket-system/lib/logger"
	"ticket-system/notification/internal/config"
)

type orderMessenger struct {
	reader  *kafka.Reader
	handler Handler
}

func NewOrderMessenger(handler Handler) Messenger {
	m := &orderMessenger{}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: config.Data.Brokers,
		Topic:   config.Data.Topics.Order,
		GroupID: config.Data.Consumers.Order,
	})

	m.reader = r
	m.handler = handler

	return m
}

func (m *orderMessenger) Close() error {
	return m.reader.Close()
}

func (m *orderMessenger) Consume(ctx context.Context, handleErr func(err error)) {
	for {
		message, err := m.reader.ReadMessage(ctx)
		if err != nil {
			handleErr(err)
			continue
		}

		err = m.handle(ctx, message)
		if err != nil {
			handleErr(err)
		}
	}
}

func (m *orderMessenger) handle(ctx context.Context, message kafka.Message) error {
	var messageType string
	messageTypeHeaderKey := string(desc.Headers_HEADER_MESSAGE_TYPE)
	for _, header := range message.Headers {
		if header.Key == messageTypeHeaderKey {
			messageType = string(header.Value)
			break
		}
	}

	if messageType == "" {
		return ErrMessageTypeHeaderNotFound
	}

	switch messageType {
	case string(desc.MessageType_ORDER_MESSAGE):
		return m.handleOrderMessage(ctx, message)
	}

	return ErrUnknownMessageType
}

func (m *orderMessenger) handleOrderMessage(ctx context.Context, message kafka.Message) error {
	orderMessage := &desc.OrderMessage{}
	err := proto.Unmarshal(message.Value, orderMessage)
	if err != nil {
		logger.Error("unable to unmarshal message", zap.Error(err))
		return err
	}

	return m.handler.HandleOrderMessage(ctx, orderMessage)
}

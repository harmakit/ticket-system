package message

import (
	"context"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"ticket-system/checkout/internal/config"
	"ticket-system/checkout/internal/model"
	desc "ticket-system/checkout/pkg/v1/message"
	"ticket-system/lib/logger"
)

type OrderMessenger interface {
	Send(ctx context.Context, order *model.Order) error
}

type orderMessenger struct {
	writer *kafka.Writer
}

func NewOrderMessenger() OrderMessenger {
	m := &orderMessenger{}
	err := m.createPartition()
	if err != nil {
		logger.Error("unable to create partition", zap.Error(err))
	}

	w := &kafka.Writer{
		Addr:         kafka.TCP(config.Data.Brokers...),
		Topic:        config.Data.Topics.Order,
		Balancer:     &kafka.RoundRobin{},
		RequiredAcks: kafka.RequireOne,
		Completion: func(messages []kafka.Message, err error) {
			if err != nil {
				logger.Error("Message not written to kafka", zap.Error(err))
				return
			}
			logger.Info("Message written to kafka", zap.Int("msg_count", len(messages)))
		},
	}
	m.writer = w

	return m
}

func (m *orderMessenger) createPartition() error {
	conn, err := kafka.DialLeader(context.Background(), "tcp", config.Data.Brokers[0], config.Data.Topics.Order, 0)
	if err != nil {
		panic(err.Error())
	}

	err = conn.Close()
	if err != nil {
		return err
	}
	return nil
}

func (m *orderMessenger) Send(ctx context.Context, o *model.Order) error {
	key := []byte(o.Id)
	message := m.getDescOrder(o)

	val, err := proto.Marshal(message)
	if err != nil {
		return err
	}

	return m.writer.WriteMessages(ctx, kafka.Message{Key: key, Value: val})
}

func (m *orderMessenger) getDescOrder(o *model.Order) *desc.Order {
	var status desc.Status
	switch o.Status {
	case model.StatusCreated:
		status = desc.Status_CREATED
	case model.StatusPaid:
		status = desc.Status_PAID
	case model.StatusFailed:
		status = desc.Status_FAILED
	case model.StatusCancelled:
		status = desc.Status_CANCELLED
	}

	return &desc.Order{
		Id:     string(o.Id),
		UserId: string(o.UserId),
		Status: status,
	}
}

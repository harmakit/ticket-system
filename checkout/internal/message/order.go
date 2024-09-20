package message

import (
	"context"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"net"
	"strconv"
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
	err := m.createTopic()
	if err != nil {
		logger.Error("unable to create partition", zap.Error(err))
	}

	w := &kafka.Writer{
		Addr:         kafka.TCP(config.Data.Brokers...),
		Topic:        config.Data.Topics.Order.Name,
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

func (m *orderMessenger) createTopic() error {
	conn, err := kafka.Dial("tcp", config.Data.Brokers[0])
	if err != nil {
		panic(err.Error())
	}

	controller, err := conn.Controller()
	if err != nil {
		panic(err.Error())
	}
	var controllerConn *kafka.Conn
	controllerConn, err = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		panic(err.Error())
	}

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             config.Data.Topics.Order.Name,
			NumPartitions:     config.Data.Topics.Order.Partitions,
			ReplicationFactor: config.Data.Topics.Order.ReplicationFactor,
		},
	}

	err = controllerConn.CreateTopics(topicConfigs...)
	if err != nil {
		panic(err.Error())
	}

	err = controllerConn.Close()
	if err != nil {
		return err
	}

	err = conn.Close()
	if err != nil {
		return err
	}

	return nil
}

func (m *orderMessenger) Send(ctx context.Context, o *model.Order) error {
	km, err := m.getOrderKafkaMessage(o)
	if err != nil {
		return err
	}
	return m.writer.WriteMessages(ctx, km)
}

func (m *orderMessenger) getDescOrder(o *model.Order) *desc.OrderMessage {
	var status desc.OrderStatus
	switch o.Status {
	case model.StatusCreated:
		status = desc.OrderStatus_CREATED
	case model.StatusPaid:
		status = desc.OrderStatus_PAID
	case model.StatusFailed:
		status = desc.OrderStatus_FAILED
	case model.StatusCancelled:
		status = desc.OrderStatus_CANCELLED
	}

	return &desc.OrderMessage{
		Id:     string(o.Id),
		UserId: string(o.UserId),
		Status: status,
	}
}

func (m *orderMessenger) getOrderKafkaMessage(o *model.Order) (kafka.Message, error) {
	var status desc.OrderStatus
	switch o.Status {
	case model.StatusCreated:
		status = desc.OrderStatus_CREATED
	case model.StatusPaid:
		status = desc.OrderStatus_PAID
	case model.StatusFailed:
		status = desc.OrderStatus_FAILED
	case model.StatusCancelled:
		status = desc.OrderStatus_CANCELLED
	}

	message := &desc.OrderMessage{
		Id:     string(o.Id),
		UserId: string(o.UserId),
		Status: status,
	}
	key := []byte(o.Id)

	val, err := proto.Marshal(message)
	if err != nil {
		return kafka.Message{}, err
	}

	typeHeader := kafka.Header{
		Key:   string(desc.Headers_HEADER_MESSAGE_TYPE),
		Value: []byte(string(desc.MessageType_ORDER_MESSAGE)),
	}

	km := kafka.Message{Key: key, Value: val, Headers: []kafka.Header{typeHeader}}

	return km, nil
}

package config

import (
	"github.com/pkg/errors"
	"os"
	"strconv"
)

type Config struct {
	DatabaseUrl string
	Port        int
	Env         string
	Services    struct {
		Booking string
		Event   string
	}
	Brokers []string
	Topics  struct {
		Order struct {
			Name              string
			Partitions        int
			ReplicationFactor int
		}
	}
}

var Data *Config

func Init() error {
	var err error
	Data = &Config{}

	Data.Port, err = strconv.Atoi(os.Getenv("SERVER_PORT"))
	if err != nil {
		return errors.WithMessage(err, "invalid port")
	}

	Data.DatabaseUrl = os.Getenv("DATABASE_URL")

	Data.Env = os.Getenv("ENV")

	Data.Services.Booking = os.Getenv("BOOKING_SERVER")
	Data.Services.Event = os.Getenv("EVENT_SERVER")

	Data.Brokers = []string{
		os.Getenv("BROKER_1"),
		os.Getenv("BROKER_2"),
		os.Getenv("BROKER_3"),
	}
	Data.Topics.Order.Name = os.Getenv("TOPIC_ORDER")
	Data.Topics.Order.Partitions, err = strconv.Atoi(os.Getenv("TOPIC_ORDER_PARTITIONS"))
	if err != nil {
		return errors.WithMessage(err, "invalid partitions")
	}
	Data.Topics.Order.ReplicationFactor, err = strconv.Atoi(os.Getenv("TOPIC_ORDER_REPLICATION_FACTOR"))
	if err != nil {
		return errors.WithMessage(err, "invalid replication factor")
	}

	err = Validate()
	if err != nil {
		return errors.WithMessage(err, "invalid config")
	}

	return nil
}

func Validate() error {
	if Data.Port < 0 || Data.Port > 65535 {
		return errors.New("invalid port")
	}

	if Data.Env != "dev" && Data.Env != "prod" {
		return errors.New("invalid env")
	}

	if Data.DatabaseUrl == "" {
		return errors.New("invalid database url")
	}

	if Data.Services.Booking == "" {
		return errors.New("invalid booking service")
	}
	if Data.Services.Event == "" {
		return errors.New("invalid booking service")
	}

	if len(Data.Brokers) == 0 {
		return errors.New("no brokers provided")
	}
	for _, broker := range Data.Brokers {
		if broker == "" {
			return errors.New("invalid broker")
		}
	}

	if Data.Topics.Order.Name == "" {
		return errors.New("invalid order topic")
	}
	if Data.Topics.Order.Partitions == 0 {
		return errors.New("invalid order partitions")
	}
	if Data.Topics.Order.ReplicationFactor == 0 {
		return errors.New("invalid order replication factor")
	}

	return nil
}

func (*Config) IsDev() bool {
	return Data.Env == "dev"
}

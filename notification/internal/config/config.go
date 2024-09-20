package config

import (
	"github.com/pkg/errors"
	"os"
)

type Config struct {
	Env     string
	Brokers []string
	Topics  struct {
		Order string
	}
	Consumers struct {
		Order string
	}
}

var Data *Config

func Init() error {
	var err error
	Data = &Config{}

	Data.Env = os.Getenv("ENV")

	Data.Brokers = []string{
		os.Getenv("BROKER_1"),
		os.Getenv("BROKER_2"),
		os.Getenv("BROKER_3"),
	}
	Data.Topics.Order = os.Getenv("TOPIC_ORDER")
	Data.Consumers.Order = os.Getenv("TOPIC_ORDER_CONSUMER_GROUP")

	err = Validate()
	if err != nil {
		return errors.WithMessage(err, "invalid config")
	}

	return nil
}

func Validate() error {

	if Data.Env != "dev" && Data.Env != "prod" {
		return errors.New("invalid env")
	}

	if len(Data.Brokers) == 0 {
		return errors.New("no brokers provided")
	}
	for _, broker := range Data.Brokers {
		if broker == "" {
			return errors.New("invalid broker")
		}
	}

	if Data.Topics.Order == "" {
		return errors.New("invalid order topic")
	}
	if Data.Consumers.Order == "" {
		return errors.New("invalid order consumer group")
	}

	return nil
}

func (*Config) IsDev() bool {
	return Data.Env == "dev"
}

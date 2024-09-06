package config

import (
	"github.com/pkg/errors"
	"os"
	"strconv"
)

type Config struct {
	Port int
	Env  string
}

var Data *Config

func Init() error {
	var err error
	Data = &Config{}

	Data.Port, err = strconv.Atoi(os.Getenv("SERVER_PORT"))
	if err != nil {
		return errors.WithMessage(err, "invalid port")
	}

	Data.Env = os.Getenv("ENV")

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

	return nil
}

func (*Config) IsDev() bool {
	return Data.Env == "dev"
}

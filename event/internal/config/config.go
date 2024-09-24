package config

import (
	"github.com/pkg/errors"
	"os"
	"strconv"
	"sync"
)

type Config struct {
	DatabaseUrl string
	Port        int
	Env         string
	Services    struct {
		Booking string
	}
}

var Data *Config
var once sync.Once

func Init() error {
	var err error
	once.Do(func() {
		Data = &Config{}

		Data.Port, err = strconv.Atoi(os.Getenv("SERVER_PORT"))
		if err != nil {
			err = errors.WithMessage(err, "invalid port")
			return
		}

		Data.DatabaseUrl = os.Getenv("DATABASE_URL")

		Data.Env = os.Getenv("ENV")

		Data.Services.Booking = os.Getenv("BOOKING_SERVER")

		err = Validate()
		if err != nil {
			err = errors.WithMessage(err, "invalid config")
			return
		}
	})

	return err
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

	return nil
}

func (*Config) IsDev() bool {
	return Data.Env == "dev"
}

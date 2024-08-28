package config

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Port int
	Env  string
}

var Data *Config

func Init() error {
	rawYML, err := os.ReadFile("config.yml")
	if err != nil {
		return errors.WithMessage(err, "unable to read config.yml")
	}

	err = yaml.Unmarshal(rawYML, &Data)
	if err != nil {
		return errors.WithMessage(err, "unable to unmarshal config.yml")
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

	return nil
}

func (*Config) IsDev() bool {
	return Data.Env == "dev"
}

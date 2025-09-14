package configs

import (
	"github.com/caarlos0/env/v6"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func LoadConfig(cfg any) error {
	if err := env.Parse(cfg); err != nil {
		return err
	}
	if err := validate.Struct(cfg); err != nil {
		return err
	}
	return nil
}

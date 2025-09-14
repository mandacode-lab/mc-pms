package configs

import (
	"strconv"
	"time"
)

type ServerConfig struct {
	Port         int           `env:"PORT" envDefault:"8080"`
	ReadTimeout  time.Duration `env:"READ_TIMEOUT" envDefault:"30s"`
	WriteTimeout time.Duration `env:"WRITE_TIMEOUT" envDefault:"30s"`
}

func (c *ServerConfig) GenerateAddress() string {
	return ":" + strconv.Itoa(c.Port)
}

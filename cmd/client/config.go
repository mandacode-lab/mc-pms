package main

import (
	"github.com/mandacode-lab/mc-pms/internal/config"
)

// Config represents the client service configuration
type Config struct {
	HTTP config.HTTP `envPrefix:"HTTP_"`
	DB   config.DB   `envPrefix:"DB_"`
}

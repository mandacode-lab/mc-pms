package main

import (
	"github.com/mandacode-com/mandacode-pms/internal/config"
)

// Config represents the admin service configuration
type Config struct {
	HTTP config.HTTP `envPrefix:"HTTP_"`
	DB   config.DB   `envPrefix:"DB_"`
}

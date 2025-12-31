package config

import "time"

// HTTP represents HTTP server configuration
type HTTP struct {
	Port            string        `env:"PORT" envDefault:"8080"`
	ReadTimeout     time.Duration `env:"READ_TIMEOUT" envDefault:"10s"`
	WriteTimeout    time.Duration `env:"WRITE_TIMEOUT" envDefault:"10s"`
	IdleTimeout     time.Duration `env:"IDLE_TIMEOUT" envDefault:"60s"`
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT" envDefault:"5s"`
	Host            string        `env:"HOST" envDefault:"0.0.0.0"`
	CORS            CORS          `envPrefix:"CORS_"`
}

// CORS represents CORS configuration
type CORS struct {
	Enabled          bool     `env:"ENABLED" envDefault:"false"`
	AllowedOrigins   []string `env:"ALLOWED_ORIGINS" envSeparator:"," envDefault:"*"`
	AllowedMethods   []string `env:"ALLOWED_METHODS" envSeparator:"," envDefault:"GET,POST,PUT,PATCH,DELETE,OPTIONS"`
	AllowedHeaders   []string `env:"ALLOWED_HEADERS" envSeparator:"," envDefault:"Origin,Content-Length,Content-Type,Authorization,Accept,X-Requested-With"`
	ExposeHeaders    []string `env:"EXPOSE_HEADERS" envSeparator:","`
	AllowCredentials bool     `env:"ALLOW_CREDENTIALS" envDefault:"false"`
	MaxAge           int      `env:"MAX_AGE" envDefault:"43200"` // 12 hours in seconds
}

// Address returns the server address
func (c HTTP) Address() string {
	return c.Host + ":" + c.Port
}

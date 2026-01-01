package config

// Postgres represents PostgreSQL database configuration
type DB struct {
	Host     string `env:"HOST" envDefault:"localhost"`
	Port     string `env:"PORT" envDefault:"5432"`
	Username string `env:"USERNAME" envDefault:"postgres"`
	Password string `env:"PASSWORD" envDefault:"postgres"`
	Database string `env:"DATABASE" envDefault:"postgres"`
	SSLMode  string `env:"SSL_MODE" envDefault:"disable"`
}

// Address returns the PostgreSQL connection string
func (c DB) Address() string {
	return "host=" + c.Host + " port=" + c.Port + " user=" + c.Username + " password=" + c.Password + " dbname=" + c.Database + " sslmode=" + c.SSLMode
}

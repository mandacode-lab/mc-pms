package configs

import "strconv"

type PostgresConfig struct {
	Host     string `env:"POSTGRES_HOST" envDefault:"localhost"`
	Port     int    `env:"POSTGRES_PORT" envDefault:"5432"`
	User     string `env:"POSTGRES_USER" envDefault:"postgres"`
	Password string `env:"POSTGRES_PASSWORD" envDefault:""`
	Database string `env:"POSTGRES_DB" envDefault:"mandacode-ssam"`
	SSLMode  string `env:"POSTGRES_SSL_MODE" envDefault:"disable" validate:"oneof=disable require verify-ca verify-full"`
}

func (c PostgresConfig) GetDSN() string {
	return "postgres://" + c.User + ":" + c.Password + "@" + c.Host + ":" + strconv.Itoa(c.Port) + "/" + c.Database + "?sslmode=" + c.SSLMode
}

package configs

type AuthConfig struct {
	Env      string         `env:"ENV" envDefault:"prod" validate:"oneof=dev staging prod"`
	Server   ServerConfig   `envPrefix:""`
	Postgres PostgresConfig `envPrefix:""`
	Redis    RedisConfig    `envPrefix:""`
}

func LoadAuthConfig() (*AuthConfig, error) {
	cfg := &AuthConfig{}
	if err := LoadConfig(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}


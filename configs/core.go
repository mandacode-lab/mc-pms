package configs

type CoreConfig struct {
	Env      string         `env:"ENV" envDefault:"prod" validate:"oneof=dev staging prod"`
	KEK      string         `env:"KEK" validate:"required,len=64,hexadecimal"` // 256-bit hex-encoded key (64 chars)
	Server   ServerConfig   `envPrefix:""`
	Postgres PostgresConfig `envPrefix:""`
	Redis    RedisConfig    `envPrefix:""`
}

func LoadCoreConfig() (*CoreConfig, error) {
	cfg := &CoreConfig{}
	if err := LoadConfig(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

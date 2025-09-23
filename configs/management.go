package configs

type ManagementConfig struct {
	Env      string         `env:"ENV" envDefault:"prod" validate:"oneof=dev staging prod"`
	Server   ServerConfig   `envPrefix:""`
	Postgres PostgresConfig `envPrefix:""`
	Redis    RedisConfig    `envPrefix:""`
	IAM      IAMConfig      `envPrefix:"IAM_"`
}

func LoadManagementConfig() (*ManagementConfig, error) {
	cfg := &ManagementConfig{}
	if err := LoadConfig(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

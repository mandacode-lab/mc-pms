package configs

type InitializeConfig struct {
	Env          string             `env:"ENV" envDefault:"prod" validate:"oneof=dev staging prod"`
	Postgres     PostgresConfig     `envPrefix:""`
	SystemClient SystemClientConfig `envPrefix:""`
}

func LoadInitializeConfig() (*InitializeConfig, error) {
	cfg := &InitializeConfig{}
	if err := LoadConfig(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
